package lambda

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/constant"
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
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	lambTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	scheTypes "github.com/aws/aws-sdk-go-v2/service/scheduler/types"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

//go:generate mockgen -destination ../../testdata/lambda_service_mock.go -package testdata -source service.go Service
type (
	LambdaService interface {
		Register(c context.Context, r *dto.ReqRegister) ([]*dto.RespRegister, error)
		Invoke(c context.Context, r *dto.ReqInvoke) (*dto.RespInvoke, error)
		List(c context.Context, isFull bool) (interface{}, error)
		Info(c context.Context, r *dto.ReqURILambda) (*dto.RespInfo, error)
		Logs(c context.Context, r *dto.ReqURILambda, upgrader *websocket.Upgrader) error
		Remove(c context.Context, r *dto.ReqURILambda) (*dto.RespRemove, error)
	}
	service struct {
		lambdaRepo repo.Lambda
		amazon     amazonx.Amazon
		oauthRepo  repo.OAuth
	}
)

var LambdaServiceImpl LambdaService

func NewLambdaService() {
	if LambdaServiceImpl == nil {
		repo.NewLambda()

		LambdaServiceImpl = &service{
			lambdaRepo: repo.LambdaRepo,
			amazon:     amazonx.Conductor,
			oauthRepo:  repo.OAuthRepo,
		}
	}
}

type toBePersistPair struct {
	Lambda    *model.Lambda
	Scheduler *model.LambdaScheduler
}

func (svc *service) Register(c context.Context, r *dto.ReqRegister) ([]*dto.RespRegister, error) {
	expression := r.Expression
	files := r.Files

	jwtOrg, _ := c.(*gin.Context).Get(constant.ClaimIss.Str())
	jwtAccount, _ := c.(*gin.Context).Get(constant.ClaimSub.Str())

	user, err := svc.oauthRepo.FindUserByOrgAcn(c, &dto.ReqOrgAcn{
		OrgName: jwtOrg.(string),
		AcnName: jwtAccount.(string),
	})
	if err != nil {
		return nil, err
	}

	maxLimit := config.GlobalConfig.Lambda.Max
	ls, err := svc.lambdaRepo.FindByAccount(c, user.ID)
	if err != nil {
		return nil, err
	}
	if len(ls)+len(files) > maxLimit {
		return nil, errorx.BadRequest(fmt.Sprintf("the number of lambdas is limited to %d", maxLimit))
	}

	toBePersist := make([]toBePersistPair, 0, len(files))
	resp := make([]*dto.RespRegister, 0, len(files))

	roleName := util.GetRoleName(c, jwtOrg.(string), jwtAccount.(string))
	roleARN, err := svc.getRoleARN(c, roleName)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		newLamResp, err := svc.registerLambda(c, file, roleARN)
		if err != nil {
			return nil, err
		}

		respItem := &dto.RespRegister{Lambda: &dto.RespLamBrief{
			Name:    *newLamResp.FunctionName,
			Arn:     *newLamResp.FunctionArn,
			Runtime: string(newLamResp.Runtime),
			Handler: *newLamResp.Handler,
			Version: *newLamResp.Version,
		}}

		tpp := toBePersistPair{
			Lambda: model.BuildLambda(
				model.WithLambdaResp(newLamResp),
				model.WithAccountID(user.ID),
			),
		}

		if strings.TrimSpace(expression) != "" {
			newSchResp, err := svc.boundScheduler(c, newLamResp, expression, r.Payload, roleARN)
			if err != nil {
				return nil, err
			}

			respItem.Scheduler = &dto.RespSchBrief{
				Arn:            *newSchResp.ScheduleArn,
				Name:           *newLamResp.FunctionName,
				BoundLambdaArn: *newLamResp.FunctionArn,
			}

			tpp.Scheduler = model.BuildScheduler(
				model.WithExpression(expression),
				model.WithSchArn(*newSchResp.ScheduleArn),
				// in binding the scheduler with Lambda, the scheduler name is from the name of Lambda.
				model.WithSchName(*newLamResp.FunctionName),
			)
		} else {
			logx.Logger.INFO(fmt.Sprintf("%s: will be triggered manually", file.Name))
		}

		toBePersist = append(toBePersist, tpp)
		resp = append(resp, respItem)
	}

	if err := svc.persistRegisterResults(c, toBePersist); err != nil {
		return nil, err
	}

	return resp, nil
}

func (svc *service) registerLambda(c context.Context, file *dto.ReqFile, roleARN string) (*lambda.CreateFunctionOutput, error) {
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
			Role:        aws.String(roleARN),
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
	inputPayload string,
	roleARN string,
) (*scheduler.CreateScheduleOutput, error) {
	event, err := util.GenEventPayload(c, inputPayload)
	if err != nil {
		return nil, err
	}
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return nil, errorx.Internal(fmt.Sprintf("failed to marshal event payload to json: %s", err.Error()))
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
			RoleArn: aws.String(roleARN),
			Input:   aws.String(string(eventJSON)),
		},
		ActionAfterCompletion:      scheTypes.ActionAfterCompletionNone,
		Description:                nil,
		ScheduleExpressionTimezone: aws.String("UTC"),
		State:                      scheTypes.ScheduleStateEnabled,
	})
	if err != nil {
		return nil, errorx.Internal(fmt.Sprintf("failed to bound scheduler: %s, err: %s", *lambdaFun.FunctionName, err.Error()))
	}

	logx.Logger.DEBUG(fmt.Sprintf("scheduler created: %s", *newSchResp.ScheduleArn))

	return newSchResp, nil
}

func (svc *service) persistRegisterResults(c context.Context, pairs []toBePersistPair) error {
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
		logx.Logger.ERROR("failed to persistRegisterResults lambda related data", map[string]interface{}{"error": err.Error()})
		return errorx.Internal(err.Error())
	}

	return nil
}

func (svc *service) Invoke(c context.Context, r *dto.ReqInvoke) (*dto.RespInvoke, error) {
	jwtAccount, _ := c.(*gin.Context).Get(constant.ClaimSub.Str())

	user, err := svc.oauthRepo.FindUserByAcn(c, jwtAccount.(string))
	if err != nil {
		return nil, err
	}

	lamb, err := svc.lambdaRepo.LambdaInfo(c, user.ID, r.Lambda)
	if err != nil {
		return nil, err
	}

	payload, err := util.GenEventPayload(c, r.Payload)
	if err != nil {
		return nil, err
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
		return nil, errorx.Internal(fmt.Sprintf("failed to invoke lambda: %s, error: %s", lamb.FunctionName, err.Error()))
	}

	// decode log result
	decodedLogResult, err := util.DecodeBase64String(c, invokeOutput.LogResult)
	if err != nil {
		return nil, errorx.Internal(fmt.Sprintf("failed to decode log result: %s", err.Error()))
	}
	invokeOutput.LogResult = &decodedLogResult

	// decode payload
	decodedPayload, err := util.DecodeBase64String(c, aws.String(string(invokeOutput.Payload)))
	if err != nil {
		return nil, errorx.Internal(fmt.Sprintf("failed to decode payload: %s", err.Error()))
	}
	invokeOutput.Payload = []byte(decodedPayload)

	return dto.BuildRespInvoke(dto.WithInvokeResp(invokeOutput)), nil
}

func (svc *service) List(c context.Context, isFull bool) (interface{}, error) {
	jwtAccount, _ := c.(*gin.Context).Get(constant.ClaimSub.Str())

	user, err := svc.oauthRepo.FindUserByAcn(c, jwtAccount.(string))
	if err != nil {
		return nil, err
	}

	lambs, err := svc.lambdaRepo.FindByAccount(c, user.ID)
	if err != nil {
		return nil, err
	}

	if isFull {
		return lambs, nil
	}

	respBrief := make([]*dto.RespInList, 0, len(lambs))
	for _, lamb := range lambs {
		respBrief = append(respBrief, &dto.RespInList{
			FunctionName: lamb.FunctionName,
			FunctionArn:  lamb.FunctionArn,
			Description:  lamb.Description,
			CreatedAt:    lamb.CreatedAt,
		})
	}

	return respBrief, nil
}

func (svc *service) Info(c context.Context, r *dto.ReqURILambda) (*dto.RespInfo, error) {
	jwtAccount, _ := c.(*gin.Context).Get(constant.ClaimSub.Str())

	user, err := svc.oauthRepo.FindUserByAcn(c, jwtAccount.(string))
	if err != nil {
		return nil, err
	}

	info, err := svc.lambdaRepo.LambdaInfo(c, user.ID, r.Lambda)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (svc *service) Logs(c context.Context, req *dto.ReqURILambda, upgrader *websocket.Upgrader) error {
	// websocket
	ctx, ok := c.(*gin.Context)
	if !ok {
		return errorx.GinContextConv()
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

func (svc *service) Remove(c context.Context, r *dto.ReqURILambda) (*dto.RespRemove, error) {
	jwtAccount, _ := c.(*gin.Context).Get(constant.ClaimSub.Str())

	user, err := svc.oauthRepo.FindUserByAcn(c, jwtAccount.(string))
	if err != nil {
		return nil, err
	}

	lamb, err := svc.lambdaRepo.LambdaInfo(c, user.ID, r.Lambda)
	if err != nil {
		return nil, err
	}

	rmvLambInput := &lambda.DeleteFunctionInput{
		FunctionName: aws.String(lamb.FunctionName),
	}
	rmvLamb, err := svc.amazon.RemoveLambda(c, rmvLambInput)
	if err != nil {
		return nil, err
	}
	logx.Logger.INFO(fmt.Sprintf(
		"lambda <%s/%s> removed\nmetadata: %v",
		lamb.FunctionName, lamb.FunctionArn,
		rmvLamb.ResultMetadata,
	))

	if lamb.Scheduler.ScheduleArn != "" {
		rmvSchInput := &scheduler.DeleteScheduleInput{
			Name: aws.String(lamb.Scheduler.ScheduleName),
		}
		rmvSch, err := svc.amazon.RemoveScheduler(c, rmvSchInput)
		if err != nil {
			return nil, err
		}
		logx.Logger.INFO(fmt.Sprintf(
			"scheduler <%s/%s> removed metadata %v",
			lamb.Scheduler.ScheduleName, lamb.Scheduler.ScheduleArn,
			rmvSch.ResultMetadata,
		))
	}

	if err := svc.lambdaRepo.DeleteLambdaTX(
		c,
		func(tx *gorm.DB) error {
			if err := tx.
				Where(map[string]interface{}{
					"function_arn": lamb.FunctionArn,
				}).
				Delete(&model.Lambda{}).
				Error; err != nil {
				return errorx.Internal(fmt.Sprintf("failed to delete lambda: %s", lamb.FunctionArn))
			}

			if err := tx.
				Where(map[string]interface{}{
					"lambda_id": lamb.ID,
				}).
				Delete(&model.LambdaScheduler{}).Error; err != nil {
				return errorx.Internal(fmt.Sprintf("failed to delete scheduler: %s", lamb.Scheduler.ScheduleArn))
			}

			return nil
		},
		&sql.TxOptions{
			Isolation: sql.LevelSerializable,
		},
	); err != nil {
		return nil, err
	}

	return &dto.RespRemove{
		Lambdas: dto.RespLamBrief{
			Name: lamb.FunctionName,
			Arn:  lamb.FunctionArn,
		},
		Scheduler: dto.RespSchBrief{
			Arn: lamb.Scheduler.ScheduleArn,
		},
	}, nil
}

func (svc *service) getRoleARN(c context.Context, roleName string) (string, error) {
	input := &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	}

	result, err := svc.amazon.GetRole(c, input)
	if err != nil {
		return "", errorx.Internal(fmt.Sprintf("get role arn error: %v", err))
	}

	return *result.Role.Arn, nil
}
