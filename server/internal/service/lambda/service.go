package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	configx "github.com/57blocks/auto-action/server/internal/config"
	pkgLog "github.com/57blocks/auto-action/server/internal/pkg/log"
	dto "github.com/57blocks/auto-action/server/internal/service/dto/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	lambTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	scheTypes "github.com/aws/aws-sdk-go-v2/service/scheduler/types"
	"github.com/pkg/errors"
)

type (
	Service interface {
		Register(c context.Context, r *http.Request) (*dto.RespRegister, error)
	}
	conductor struct{}
)

var (
	Conductor Service

	// TODO: init them for once
	awsConfig    aws.Config
	lambdaClient *lambda.Client
)

func init() {
	if Conductor == nil {
		Conductor = &conductor{}
	}
}

func (cd *conductor) Register(c context.Context, r *http.Request) (*dto.RespRegister, error) {
	var err error

	fileHeaders := r.MultipartForm.File
	expression := r.Form.Get("expression")

	awsConfig, err = config.LoadDefaultConfig(
		c,
		config.WithRegion(configx.Global.Region),
		config.WithSharedConfigProfile("iamp3ngf3i"), // TODO: only for local
	)
	if err != nil {
		pkgLog.Logger.ERROR(fmt.Sprintf("failed to load AWS config: %s", err.Error()))
		return nil, errors.New(err.Error())
	}

	resp := make([]*lambda.CreateFunctionOutput, 0, len(fileHeaders))

	for _, fhs := range r.MultipartForm.File {
		fh := fhs[0]

		cfo, err := registerLambda(c, fh)
		if err != nil {
			return nil, err
		}
		resp = append(resp, cfo)

		if strings.TrimSpace(expression) != "" {
			pkgLog.Logger.DEBUG(fmt.Sprintf("scheduler expression found: %s", expression))
			cso, err := boundScheduler(c, cfo, expression)
			if err != nil {
				return nil, err
			}
			pkgLog.Logger.DEBUG(fmt.Sprintf("scheduler created: %s", *cso.ScheduleArn))
		} else {
			pkgLog.Logger.DEBUG("no expression found, will be triggered manually")
		}

	}

	return &dto.RespRegister{
		CFOs: resp,
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

	lambdaClient = lambda.NewFromConfig(awsConfig)

	// register lambda
	lambdaFun, err := lambdaClient.CreateFunction(
		c,
		&lambda.CreateFunctionInput{
			Code: &lambTypes.FunctionCode{
				ZipFile: fileBytes,
			},
			FunctionName: aws.String(splits[0]),
			// TODO: put into env when the infrastructure is ready, the same as `VpcConfig` below.
			Environment: &lambTypes.Environment{Variables: map[string]string{"ENV_REGION": "us-east-2"}},
			// This execution role has full access of CloudWatch and Lambda execution access.
			Role:        aws.String("arn:aws:iam::123340007534:role/LambdaExecutionRole"),
			Runtime:     lambTypes.RuntimeNodejs20x,
			Description: nil,
			Handler:     aws.String(fmt.Sprintf("%s.handler", splits[0])),
			PackageType: lambTypes.PackageTypeZip,
			Publish:     false,
			VpcConfig: &lambTypes.VpcConfig{
				Ipv6AllowedForDualStack: aws.Bool(false),
				SecurityGroupIds:        []string{"sg-0b77e2d29b5bca26a"}, // default SG of the VPC
				// now below all are public subnets
				SubnetIds: []string{
					"subnet-05e584ba2ffb30d2d",
					"subnet-0ce2404213f76db94",
					"subnet-0ecc12c551120284b",
				},
			},
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
	schClient := scheduler.NewFromConfig(awsConfig)

	event, err := buildLambdaEvent(c)
	if err != nil {
		return nil, err
	}
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	lambScheduler, err := schClient.CreateSchedule(c, &scheduler.CreateScheduleInput{
		FlexibleTimeWindow: &scheTypes.FlexibleTimeWindow{
			Mode: scheTypes.FlexibleTimeWindowModeOff,
		},
		Name:               aws.String(fmt.Sprintf("scheduler-%s", *lambdaFun.FunctionName)),
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

	return lambScheduler, nil
}

//// TODO: remove the zip file function when the input file is a zip file already
//const (
//	TempZip = "handler.zip"
//)
//func zipFile(fh *multipart.FileHeader, target string) error {
//	// Open the uploaded file
//	file, err := fh.Open()
//	if err != nil {
//		return err
//	}
//	defer file.Close()
//
//	// Create the ZIP file
//	zipFile, err := os.Create(target)
//	if err != nil {
//		return err
//	}
//	defer zipFile.Close()
//
//	// Create a new ZIP writer
//	zipWriter := zip.NewWriter(zipFile)
//	defer zipWriter.Close()
//
//	// Create a file entry in the ZIP archive
//	zipEntryWriter, err := zipWriter.Create(fh.Filename)
//	if err != nil {
//		return err
//	}
//
//	// Copy the content of the uploaded file into the ZIP archive
//	_, err = io.Copy(zipEntryWriter, file)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
