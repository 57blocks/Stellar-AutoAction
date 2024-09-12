package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/model"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/57blocks/auto-action/server/internal/pkg/util"
	"github.com/57blocks/auto-action/server/internal/repo"
	"github.com/57blocks/auto-action/server/internal/third-party/amazonx"
	"github.com/57blocks/auto-action/server/internal/third-party/logx"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	lambTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	scheTypes "github.com/aws/aws-sdk-go-v2/service/scheduler/types"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type (
	Service interface {
		Register(c context.Context, r *dto.ReqRegister) (*dto.RespRegister, error)
		Invoke(c context.Context, r *dto.ReqInvoke) (*dto.RespInvoke, error)
		Info(c context.Context, r *dto.ReqInfo) (*dto.RespInfo, error)
		Logs(c context.Context, r *dto.ReqLogs) error
	}
	service struct {
		lambdaRepo repo.Lambda
		amazon     amazonx.Amazon
		oauthRepo  repo.OAuth
	}
)

var ServiceImpl Service

func NewLambdaService() {
	if ServiceImpl == nil {
		repo.NewLambda()

		ServiceImpl = &service{
			lambdaRepo: repo.LambdaRepo,
			amazon:     amazonx.Conductor,
			oauthRepo:  repo.OAuthRepo,
		}
	}
}

type toPersistPair struct {
	Lambda    *model.Lambda
	Scheduler *model.LambdaScheduler
}

func (svc *service) Register(c context.Context, r *dto.ReqRegister) (*dto.RespRegister, error) {
	expression := r.Expression
	files := r.Files

	jwtOrg, _ := c.(*gin.Context).Get("jwt_organization")
	jwtAccount, _ := c.(*gin.Context).Get("jwt_account")

	user, err := svc.oauthRepo.FindUserByOrgAcn(c, &dto.ReqOrgAcn{
		OrgName: jwtOrg.(string),
		AcnName: jwtAccount.(string),
	})
	if err != nil {
		return nil, err
	}

	max := config.GlobalConfig.Lambda.Max
	ls, err := svc.lambdaRepo.FindByAccount(c, user.ID)
	if err != nil {
		return nil, err
	}
	if len(ls)+len(files) > max {
		return nil, errorx.BadRequest(fmt.Sprintf("the number of lambdas is limited to %d", max))
	}

	// db persistence
	toPersists := make([]toPersistPair, 0, len(files))

	// brief response
	lsBrief := make([]dto.RespLamBrief, 0, len(files))
	ssBrief := make([]dto.RespSchBrief, 0, len(files))

	for _, file := range files {
		newLamResp, err := svc.registerLambda(c, file)
		if err != nil {
			return nil, err
		}
		lsBrief = append(lsBrief, dto.RespLamBrief{
			AccountId: user.ID,
			Name:      *newLamResp.FunctionName,
			Arn:       *newLamResp.FunctionArn,
			Runtime:   string(newLamResp.Runtime),
			Handler:   *newLamResp.Handler,
			Version:   *newLamResp.Version,
		})

		tpp := toPersistPair{
			Lambda: model.BuildLambda(
				model.WithLambdaResp(newLamResp),
			),
		}

		if strings.TrimSpace(expression) != "" {
			newSchResp, err := svc.boundScheduler(c, newLamResp, expression)
			if err != nil {
				return nil, err
			}

			ssBrief = append(ssBrief, dto.RespSchBrief{
				Arn:            *newLamResp.FunctionArn,
				BoundLambdaArn: *newSchResp.ScheduleArn,
			})
			tpp.Scheduler = model.BuildScheduler(
				model.WithExpression(expression),
				model.WithSchArn(*newSchResp.ScheduleArn),
			)
		} else {
			logx.Logger.INFO(fmt.Sprintf("%s: will be triggered manually", file.Name))
		}

		toPersists = append(toPersists, tpp)
	}

	go svc.persist(c, toPersists)

	return &dto.RespRegister{
		Lambdas:    lsBrief,
		Schedulers: ssBrief,
	}, nil
}

func (svc *service) registerLambda(c context.Context, file *dto.ReqFile) (*lambda.CreateFunctionOutput, error) {
	splits := strings.Split(file.Name, ".")
	fileName := splits[0]

	// register lambda
	lambdaFun, err := svc.amazon.RegisterLambda(
		c,
		&lambda.CreateFunctionInput{
			Code: &lambTypes.FunctionCode{
				ZipFile: file.Bytes,
			},
			FunctionName: aws.String(util.GenLambdaFuncName(c, fileName)),
			// TODO: put into env when the infrastructure is ready, the same as `VpcConfig` below.
			Environment: &lambTypes.Environment{Variables: map[string]string{"ENV_AWS_REGION": "us-east-2"}},
			// This execution role has full access of CloudWatch and Lambda execution access.
			Role:        aws.String("arn:aws:iam::123340007534:role/LambdaExecutionRole"),
			Runtime:     lambTypes.RuntimeNodejs20x,
			Timeout:     aws.Int32(30),
			Description: nil,
			Handler:     aws.String(fmt.Sprintf("%s.handler", fileName)),
			PackageType: lambTypes.PackageTypeZip,
			Publish:     false,
		},
		func(opt *lambda.Options) {},
	)
	if err != nil {
		return nil, errorx.Internal(fmt.Sprintf("failed to register lambda: %s, err: %s", fileName, err.Error()))
	}

	return lambdaFun, nil
}

func (svc *service) boundScheduler(
	c context.Context,
	lambdaFun *lambda.CreateFunctionOutput,
	expression string,
) (*scheduler.CreateScheduleOutput, error) {
	logx.Logger.DEBUG(fmt.Sprintf("scheduler expression found: %s", expression))

	event, err := util.GenEventPayload(c)
	if err != nil {
		return nil, err
	}
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	newSchResp, err := svc.amazon.BoundScheduler(c, &scheduler.CreateScheduleInput{
		FlexibleTimeWindow: &scheTypes.FlexibleTimeWindow{
			Mode: scheTypes.FlexibleTimeWindowModeOff,
		},
		Name:               aws.String(*lambdaFun.FunctionName),
		ScheduleExpression: aws.String(expression), // rate(1 minutes)/cron(...)
		Target: &scheTypes.Target{
			Arn: lambdaFun.FunctionArn,
			// This role has the Lambda invoke access to all Lambda functions in current AWS account.
			RoleArn: aws.String("arn:aws:iam::123340007534:role/service-role/Amazon_EventBridge_Scheduler_LAMBDA_25a70bed22"),
			Input:   aws.String(string(eventJSON)),
		},
		ActionAfterCompletion:      scheTypes.ActionAfterCompletionNone,
		Description:                nil,
		ScheduleExpressionTimezone: aws.String("UTC"),
		State:                      scheTypes.ScheduleStateEnabled,
	})
	if err != nil {
		return nil, errorx.Internal(fmt.Sprintf("failed to bind scheduler: %s, err: %s", *lambdaFun.FunctionName, err.Error()))
	}

	logx.Logger.DEBUG(fmt.Sprintf("scheduler created: %s", *newSchResp.ScheduleArn))

	return newSchResp, nil
}

func (svc *service) persist(c context.Context, pairs []toPersistPair) {
	if err := svc.lambdaRepo.PersistRegResult(c, func(tx *gorm.DB) error {
		for _, pair := range pairs {
			newLambda := tx.Table("lambda").Create(pair.Lambda)
			if err := newLambda.Error; err != nil {
				return errorx.Internal(fmt.Sprintf("failed to create lambda: %s", pair.Lambda.FunctionArn))
			}

			if pair.Scheduler == nil {
				continue
			}

			pair.Scheduler.LambdaID = pair.Lambda.ID

			if err := tx.Table("lambda_scheduler").Create(pair.Scheduler).Error; err != nil {
				return errorx.NotFound(fmt.Sprintf("failed to create lambda scheduler: %s",
					pair.Scheduler.ScheduleArn))
			}
		}

		return nil
	}); err != nil {
		logx.Logger.ERROR("failed to persist lambda related data", map[string]interface{}{"error": err.Error()})
	}
}

func (svc *service) Invoke(c context.Context, r *dto.ReqInvoke) (*dto.RespInvoke, error) {
	lamb, err := svc.lambdaRepo.FindByNameOrARN(c, r.Lambda)

	// generate merged payload with orgSecretKey key
	stdPayload, err := util.GenEventPayload(c)
	if err != nil {
		return nil, err
	}
	payload := map[string]interface{}{
		"organization": stdPayload.Organization,
		"account":      stdPayload.Account,
	}

	if len([]byte(r.Payload)) > 0 {
		var inputPayload map[string]interface{}
		if err := json.Unmarshal([]byte(r.Payload), &inputPayload); err != nil {
			return nil, errorx.Internal(fmt.Sprintf("failed to unmarshal input payload: %s", err.Error()))
		}

		payload = util.MergeMaps(payload, inputPayload)
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, errorx.Internal(fmt.Sprintf("failed to marshal payload: %s", err.Error()))
	}

	// invoke
	invokeOutput, err := svc.amazon.InvokeLambda(c, &lambda.InvokeInput{
		FunctionName: aws.String(lamb.FunctionName),
		LogType:      lambTypes.LogTypeTail,
		Payload:      payloadBytes,
	})
	if err != nil {
		return nil, errorx.Internal(fmt.Sprintf("failed to invoke lambda: %s", lamb.FunctionName))
	}

	return dto.BuildRespInvoke(dto.WithInvokeResp(invokeOutput)), nil
}

func (svc *service) Info(c context.Context, r *dto.ReqInfo) (*dto.RespInfo, error) {
	info, err := svc.lambdaRepo.LambdaInfo(c, r)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (svc *service) Logs(c context.Context, req *dto.ReqLogs) error {
	// websocket
	ctx, ok := c.(*gin.Context)
	if !ok {
		return errorx.GinContextConv()
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	wsConn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return errorx.Internal(fmt.Sprintf("failed to upgrade websocket: %s", err.Error()))
	}
	defer wsConn.Close()

	logGroupName := "/aws/lambda/" + util.GenLambdaFuncName(c, req.Lambda)

	describeInput := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: &logGroupName,
		OrderBy:      types.OrderByLastEventTime,
		Descending:   aws.Bool(true),
	}

	describeOutput, err := svc.amazon.DescribeLogStreams(c, describeInput)
	if err != nil {
		return errorx.Internal(fmt.Sprintf("failed to describe log streams: %s", err.Error()))
	}

	if len(describeOutput.LogStreams) == 0 {
		return errorx.NotFound("no log streams found")
	}

	logStreamName := describeOutput.LogStreams[0].LogStreamName

	input := &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  &logGroupName,
		LogStreamName: logStreamName,
	}

	var nextToken *string

	for {
		if nextToken != nil {
			input.NextToken = nextToken
		}

		output, err := svc.amazon.GetLogEvents(c, input)
		if err != nil {
			return errorx.Internal(fmt.Sprintf("failed to get log events: %s", err.Error()))
		}

		for _, event := range output.Events {
			if err := wsConn.WriteMessage(websocket.TextMessage, []byte(*event.Message)); err != nil {
				return errorx.Internal(fmt.Sprintf("failed to write message: %s", err.Error()))
			}
		}

		nextToken = output.NextForwardToken

		if nextToken == nil {
			time.Sleep(30 * time.Second)
		}
	}
}
