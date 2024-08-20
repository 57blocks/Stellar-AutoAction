package lambda

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/57blocks/auto-action/server/internal/config"
	pkgLog "github.com/57blocks/auto-action/server/internal/pkg/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.Global.Region),
	})
	if err != nil {
		pkgLog.Logger.ERROR(fmt.Sprintf("failed to create AWS session: %s", err.Error()))
		ctx.AbortWithError(
			http.StatusInternalServerError,
			err,
		)
		return
	}

	// Read Lambda function code from local file
	functionCode, err := os.ReadFile("path/to/your/lambda/function.zip")
	if err != nil {
		pkgLog.Logger.ERROR(fmt.Sprintf("failed to read Lambda function code: %s", err.Error()))
		ctx.AbortWithError(
			http.StatusInternalServerError,
			err,
		)
		return
	}

	// Create IAM role for Lambda
	iamSvc := iam.New(sess)
	roleName := "lambda-ex"
	role, err := iamSvc.CreateRole(&iam.CreateRoleInput{
		RoleName: aws.String(roleName),
		AssumeRolePolicyDocument: aws.String(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {
						"Service": "lambda.amazonaws.com"
					},
					"Action": "sts:AssumeRole"
				}
			]
		}`),
	})
	if err != nil {
		pkgLog.Logger.ERROR(fmt.Sprintf("failed to create IAM role:: %s", err.Error()))
		ctx.AbortWithError(
			http.StatusInternalServerError,
			err,
		)
		return
	}

	_, err = iamSvc.AttachRolePolicy(&iam.AttachRolePolicyInput{
		RoleName:  aws.String(roleName),
		PolicyArn: aws.String("arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"),
	})
	if err != nil {
		pkgLog.Logger.ERROR(fmt.Sprintf("failed to attach policy to IAM role: %w", err.Error()))
		ctx.AbortWithError(
			http.StatusInternalServerError,
			err,
		)
		return
	}

	// Create Lambda function
	lambdaSvc := lambda.New(sess)
	_, err = lambdaSvc.CreateFunction(&lambda.CreateFunctionInput{
		FunctionName: aws.String("my-function"),
		Runtime:      aws.String("go1.x"),
		Role:         role.Role.Arn,
		Handler:      aws.String("main"),
		Code: &lambda.FunctionCode{
			ZipFile: functionCode,
		},
	})
	if err != nil {
		pkgLog.Logger.ERROR(fmt.Sprintf("failed to create Lambda function: %w", err.Error()))
		ctx.AbortWithError(
			http.StatusInternalServerError,
			err,
		)
		return
	}

	log.Println("Successfully uploaded Lambda function")

	ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
}
