package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	configx "github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/db"
	"github.com/57blocks/auto-action/server/internal/model"
	pkgLog "github.com/57blocks/auto-action/server/internal/pkg/log"
	dto "github.com/57blocks/auto-action/server/internal/service/dto/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	lambTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	scheTypes "github.com/aws/aws-sdk-go-v2/service/scheduler/types"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type (
	Service interface {
		Register(c context.Context, r *http.Request) (*dto.RespRegister, error)
		Invoke(c context.Context, r *dto.ReqTrigger) (*dto.RespTrigger, error)
		Info(c context.Context, r *dto.ReqInfo) (*dto.RespInfo, error)
		Logs(c context.Context, r *dto.ReqLogs) error
	}
	conductor struct{}
)

var (
	Conductor Service

	// TODO: init them for once
	awsConfig    aws.Config
	lambdaClient *lambda.Client
	cwClient     *cloudwatchlogs.Client
)

func init() {
	if Conductor == nil {
		Conductor = &conductor{}
	}
}

type toPersistPair struct {
	Lambda    *model.Lambda
	Scheduler *model.LambdaScheduler
}

func (cd *conductor) Register(c context.Context, r *http.Request) (*dto.RespRegister, error) {
	var err error

	awsConfig, err = config.LoadDefaultConfig(
		c,
		config.WithRegion(configx.Global.Region),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load AWS config")
	}

	fileHeaders := r.MultipartForm.File
	expression := r.Form.Get("expression")

	// db persistence
	toPersists := make([]toPersistPair, 0, len(fileHeaders))

	// brief response
	lsBrief := make([]dto.RespLamBrief, 0, len(fileHeaders))
	ssBrief := make([]dto.RespSchBrief, 0, len(fileHeaders))

	for _, fhs := range r.MultipartForm.File {
		fh := fhs[0]

		newLamResp, err := registerLambda(c, fh)
		if err != nil {
			return nil, err
		}
		lsBrief = append(lsBrief, dto.RespLamBrief{
			Name:    *newLamResp.FunctionName,
			Arn:     *newLamResp.FunctionArn,
			Runtime: string(newLamResp.Runtime),
			Handler: *newLamResp.Handler,
			Version: *newLamResp.Version,
		})

		tpp := toPersistPair{
			Lambda: model.BuildLambda(
				model.WithLambdaResp(newLamResp),
			),
		}

		if strings.TrimSpace(expression) != "" {
			newSchResp, err := boundScheduler(c, newLamResp, expression)
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
			pkgLog.Logger.DEBUG("no expression found, will be triggered manually")
		}

		toPersists = append(toPersists, tpp)
	}

	go persist(c, toPersists)

	return &dto.RespRegister{
		Lambdas:    lsBrief,
		Schedulers: ssBrief,
	}, nil
}

func registerLambda(c context.Context, fh *multipart.FileHeader) (*lambda.CreateFunctionOutput, error) {
	file, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read zip handler file")
	}

	splits := strings.Split(fh.Filename, ".")
	fileName := splits[0]

	lambdaClient = lambda.NewFromConfig(awsConfig)

	// register lambda
	lambdaFun, err := lambdaClient.CreateFunction(
		c,
		&lambda.CreateFunctionInput{
			Code: &lambTypes.FunctionCode{
				ZipFile: fileBytes,
			},
			FunctionName: aws.String(genLambdaFuncName(c, fileName)),
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
		errMsg := fmt.Sprintf("failed to register lambda: %s, err: %s\n", fh.Filename, err.Error())
		pkgLog.Logger.ERROR(errMsg)
		return nil, errors.Wrap(err, errMsg)
	}

	return lambdaFun, nil
}

func boundScheduler(
	c context.Context,
	lambdaFun *lambda.CreateFunctionOutput,
	expression string,
) (*scheduler.CreateScheduleOutput, error) {
	pkgLog.Logger.DEBUG(fmt.Sprintf("scheduler expression found: %s", expression))

	schClient := scheduler.NewFromConfig(awsConfig)

	event, err := genEventPayload(c)
	if err != nil {
		return nil, err
	}
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	newSchResp, err := schClient.CreateSchedule(c, &scheduler.CreateScheduleInput{
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
		errMsg := fmt.Sprintf("failed to bound lambda with scheduler: %s\n", err.Error())
		pkgLog.Logger.ERROR(errMsg)
		return nil, errors.New(errMsg)
	}

	pkgLog.Logger.DEBUG(fmt.Sprintf("scheduler created: %s", *newSchResp.ScheduleArn))

	return newSchResp, nil
}

func persist(c context.Context, pairs []toPersistPair) {
	if err := db.Conn(c).Transaction(func(tx *gorm.DB) error {
		for _, pair := range pairs {
			//l := pair.Lambda
			newLambda := tx.Table("lambda").Create(pair.Lambda)
			if err := newLambda.Error; err != nil {
				return errors.Wrap(err, fmt.Sprintf("failed to create lambda: %s", pair.Lambda.FunctionArn))
			}

			if pair.Scheduler == nil {
				continue
			}

			//s := pair.Scheduler
			pair.Scheduler.LambdaID = pair.Lambda.ID

			if err := tx.Table("lambda_scheduler").Create(pair.Scheduler).Error; err != nil {
				return errors.Wrap(err, fmt.Sprintf("failed to create lambda: %s", pair.Scheduler.ScheduleArn))
			}
		}

		return nil
	}); err != nil {
		pkgLog.Logger.ERROR(err.Error())
	}
}

func (cd *conductor) Invoke(c context.Context, r *dto.ReqTrigger) (*dto.RespTrigger, error) {
	var err error

	awsConfig, err = config.LoadDefaultConfig(
		c,
		config.WithRegion(configx.Global.Region),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load AWS config")
	}

	lambdaClient = lambda.NewFromConfig(awsConfig)

	l := new(model.Lambda)

	if err := db.Conn(c).Table("lambda").
		Where(map[string]interface{}{
			"function_arn": r.Lambda,
		}).
		Or(map[string]interface{}{
			"function_name": genLambdaFuncName(c, r.Lambda),
		}).
		First(l).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(fmt.Sprintf("none lambda found by: %s\n", r.Lambda))
		}

		return nil, errors.Wrap(err, err.Error())
	}

	// generate merged payload with orgSecretKey key
	stdPayload, err := genEventPayload(c)
	if err != nil {
		return nil, err
	}

	var inputPayload map[string]interface{}
	if err := json.Unmarshal([]byte(r.Payload), &inputPayload); err != nil {
		return nil, err
	}
	inputPayload["organization"] = stdPayload.Organization
	inputPayload["account"] = stdPayload.Account
	payload, err := json.Marshal(inputPayload)
	if err != nil {
		return nil, err
	}

	// invoke
	invokeOutput, err := lambdaClient.Invoke(c, &lambda.InvokeInput{
		FunctionName: aws.String(l.FunctionName),
		LogType:      lambTypes.LogTypeTail,
		Payload:      payload,
	})
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to invoke lambda: %s", l.FunctionName))
	}

	return dto.BuildRespTrigger(dto.WithTriggerResp(invokeOutput)), nil
}

func (cd *conductor) Info(c context.Context, r *dto.ReqInfo) (*dto.RespInfo, error) {
	resp := new(dto.RespInfo)

	if err := db.Conn(c).Table("lambda").
		Preload("VPC", func(db *gorm.DB) *gorm.DB {
			return db.Table("vpc")
		}).
		Preload("Schedulers", func(db *gorm.DB) *gorm.DB {
			return db.Table("lambda_scheduler")
		}).
		Where(map[string]interface{}{
			"function_arn": r.Lambda,
		}).
		Or(map[string]interface{}{
			"function_name": r.Lambda,
		}).
		First(resp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(fmt.Sprintf("none lambda found by: %s\n", r.Lambda))
		}

		return nil, errors.Wrap(err, err.Error())
	}

	return resp, nil
}

func (cd *conductor) Logs(c context.Context, req *dto.ReqLogs) error {
	var err error

	awsConfig, err = config.LoadDefaultConfig(
		c,
		config.WithRegion(configx.Global.Region),
	)
	if err != nil {
		return errors.Wrap(err, "failed to load AWS config")
	}

	cwClient = cloudwatchlogs.NewFromConfig(awsConfig)

	// websocket
	ctx, ok := c.(*gin.Context)
	if !ok {
		return errors.New("convert context.Context to gin.Context failed")
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
		return errors.Wrap(err, "failed to set websocket upgrade")
	}
	defer wsConn.Close()

	// TODO: query Lambda name

	logGroupName := "/aws/lambda/" + genLambdaFuncName(c, req.Lambda)

	describeInput := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: &logGroupName,
		OrderBy:      types.OrderByLastEventTime,
		Descending:   aws.Bool(true),
	}

	describeOutput, err := cwClient.DescribeLogStreams(c, describeInput)
	if err != nil {
		return errors.Wrap(err, "failed to describe log streams")
	}

	if len(describeOutput.LogStreams) == 0 {
		return errors.New("no log streams found")
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

		output, err := cwClient.GetLogEvents(c, input)
		if err != nil {
			return errors.Wrap(err, "failed to get log events")
		}

		for _, event := range output.Events {
			if err := wsConn.WriteMessage(websocket.TextMessage, []byte(*event.Message)); err != nil {
				return errors.Wrap(err, "failed to write message")
			}
		}

		nextToken = output.NextForwardToken

		if nextToken == nil {
			time.Sleep(30 * time.Second)
		}
	}
}
