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

	"github.com/57blocks/auto-action/server/internal/config"
	pkgLog "github.com/57blocks/auto-action/server/internal/pkg/log"
	dto "github.com/57blocks/auto-action/server/internal/service/dto/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/pkg/errors"
)

type (
	Service interface {
		Register(c context.Context, r *http.Request) (*dto.RespRegister, error)
	}
	ServiceConductor struct{}
)

var Conductor Service

func init() {
	if Conductor == nil {
		Conductor = &ServiceConductor{}
	}
}

const (
	TempZip = "handler.zip"
)

var (
	cfg       aws.Config
	awsLambda *lambda.Client
)

func (sc *ServiceConductor) Register(c context.Context, r *http.Request) (*dto.RespRegister, error) {
	var err error

	cfg, err = awsConfig.LoadDefaultConfig(
		c,
		awsConfig.WithRegion(config.Global.Region),
		awsConfig.WithSharedConfigProfile("iamp3ngf3i"), // TODO: only for local
	)
	if err != nil {
		pkgLog.Logger.ERROR(fmt.Sprintf("failed to load AWS config: %s", err.Error()))
		return nil, errors.New(err.Error())
	}

	awsLambda = lambda.NewFromConfig(cfg)

	resp := make([]*lambda.CreateFunctionOutput, 0)

	for _, fhs := range r.MultipartForm.File {
		fh := fhs[0]

		cfo, err := registerLambda(c, fh)
		if err != nil {
			return nil, err
		}
		resp = append(resp, cfo)
	}

	return &dto.RespRegister{
		CFOs: resp,
	}, nil
}

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

func registerLambda(c context.Context, fh *multipart.FileHeader) (*lambda.CreateFunctionOutput, error) {
	err := zipFile(fh, TempZip)
	if err != nil {
		return nil, err
	}
	bs, _ := os.ReadFile(TempZip)
	defer os.Remove(TempZip)

	strs := strings.Split(fh.Filename, ".")

	// create
	function, err := awsLambda.CreateFunction(
		c,
		&lambda.CreateFunctionInput{
			Code: &types.FunctionCode{
				ZipFile: bs,
			},
			FunctionName:         aws.String(strs[0]),
			Role:                 aws.String("arn:aws:iam::123340007534:role/service-role/autoaction-role-wl50ldot"),
			Runtime:              types.RuntimeNodejs20x,
			Architectures:        nil,
			CodeSigningConfigArn: nil,
			DeadLetterConfig:     nil,
			Description:          nil,
			//Environment:          &types.Environment{Variables: map[string]string{"AWS_REGION": "us-east-2"}},
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
			VpcConfig: &types.VpcConfig{
				Ipv6AllowedForDualStack: aws.Bool(false),
				SecurityGroupIds:        []string{"sg-063f43919a7309669"},
				SubnetIds: []string{
					"subnet-05e584ba2ffb30d2d",
					"subnet-0ce2404213f76db94",
					"subnet-0ecc12c551120284b"},
			},
		},
		func(opt *lambda.Options) {},
	)
	if err != nil {
		pkgLog.Logger.ERROR(fmt.Sprintf("s3 upload failed: %s, err: %s\n", fh.Filename, err.Error()))
		return nil, err
	}

	return function, nil
}
