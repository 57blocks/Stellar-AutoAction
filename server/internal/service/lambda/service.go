package lambda

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
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
	ServiceConductor struct{}
)

var (
	Conductor Service

	// TODO: init them for once
	awsConfig    aws.Config
	lambdaClient *lambda.Client
)

func init() {
	if Conductor == nil {
		Conductor = &ServiceConductor{}
	}
}

const (
	TempZip = "handler.zip"
)

func (sc *ServiceConductor) Register(c context.Context, r *http.Request) (*dto.RespRegister, error) {
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
		}

		pkgLog.Logger.DEBUG("no expression found, will be triggered manually")
	}

	return &dto.RespRegister{
		CFOs: resp,
	}, nil
}

func registerLambda(c context.Context, fh *multipart.FileHeader) (*lambda.CreateFunctionOutput, error) {
	err := zipFile(fh, TempZip)
	if err != nil {
		return nil, err
	}
	bs, _ := os.ReadFile(TempZip)
	defer os.Remove(TempZip)

	splits := strings.Split(fh.Filename, ".")

	lambdaClient = lambda.NewFromConfig(awsConfig)

	// register lambda
	lambdaFun, err := lambdaClient.CreateFunction(
		c,
		&lambda.CreateFunctionInput{
			Code: &lambTypes.FunctionCode{
				ZipFile: bs,
			},
			FunctionName:         aws.String(splits[0]),
			Role:                 aws.String("arn:aws:iam::123340007534:role/service-role/autoaction-role-wl50ldot"),
			Runtime:              lambTypes.RuntimeNodejs20x,
			Architectures:        nil,
			CodeSigningConfigArn: nil,
			DeadLetterConfig:     nil,
			Description:          nil,
			//Environment:          &lambTypes.Environment{Variables: map[string]string{"AWS_REGION": "us-east-2"}},
			EphemeralStorage:  nil,
			FileSystemConfigs: nil,
			Handler:           aws.String("handler"),
			ImageConfig:       nil,
			KMSKeyArn:         nil,
			Layers:            nil,
			LoggingConfig:     nil,
			MemorySize:        nil,
			PackageType:       "",
			Publish:           false,
			SnapStart:         nil,
			Tags:              nil,
			Timeout:           nil,
			TracingConfig:     nil,
			VpcConfig: &lambTypes.VpcConfig{ // TODO: put into env in ECS
				Ipv6AllowedForDualStack: aws.Bool(false),
				SecurityGroupIds:        []string{"sg-063f43919a7309669"},
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
		return nil, errors.New(errMsg)
	}

	return lambdaFun, nil
}

func boundScheduler(c context.Context, lambdaFun *lambda.CreateFunctionOutput, expression string) (*scheduler.CreateScheduleOutput, error) {
	schClient := scheduler.NewFromConfig(awsConfig)

	lambScheduler, err := schClient.CreateSchedule(c, &scheduler.CreateScheduleInput{
		FlexibleTimeWindow: &scheTypes.FlexibleTimeWindow{
			Mode: scheTypes.FlexibleTimeWindowModeOff,
		},
		Name:               aws.String(fmt.Sprintf("Scheduler-%s", *lambdaFun.FunctionName)),
		ScheduleExpression: aws.String("rate(1 minutes)"),
		Target: &scheTypes.Target{
			Arn:                         lambdaFun.FunctionArn,
			RoleArn:                     aws.String("arn:aws:iam::123340007534:role/service-role/Amazon_EventBridge_Scheduler_LAMBDA_25a70bed22"),
			DeadLetterConfig:            nil,
			EcsParameters:               nil,
			EventBridgeParameters:       nil,
			Input:                       nil,
			KinesisParameters:           nil,
			RetryPolicy:                 nil,
			SageMakerPipelineParameters: nil,
			SqsParameters:               nil,
		},
		ActionAfterCompletion:      scheTypes.ActionAfterCompletionNone,
		ClientToken:                nil,
		Description:                nil,
		EndDate:                    nil,
		GroupName:                  nil,
		KmsKeyArn:                  nil,
		ScheduleExpressionTimezone: aws.String("UTC"),
		StartDate:                  nil,
		State:                      scheTypes.ScheduleStateEnabled,
	})
	if err != nil {
		errMsg := fmt.Sprintf("failed to bound lambda with scheduler: %s\n", err.Error())
		pkgLog.Logger.ERROR(errMsg)
		return nil, errors.New(errMsg)
	}

	return lambScheduler, nil
}

// TODO: remove the zip file function when the input file is a zip file already
func zipFile(fh *multipart.FileHeader, target string) error {
	// Open the uploaded file
	file, err := fh.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	// Create the ZIP file
	zipFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// Create a new ZIP writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Create a file entry in the ZIP archive
	zipEntryWriter, err := zipWriter.Create(fh.Filename)
	if err != nil {
		return err
	}

	// Copy the content of the uploaded file into the ZIP archive
	_, err = io.Copy(zipEntryWriter, file)
	if err != nil {
		return err
	}

	return nil
}
