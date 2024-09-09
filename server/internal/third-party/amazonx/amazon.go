package amazonx

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type (
	Amazon interface {
	}

	amazon struct {
		amazonConfig         aws.Config
		secretManagerClient  *secretsmanager.Client
		lambdaClient         *lambda.Client
		schedulerClient      *scheduler.Client
		cloudWatchLogsClient *cloudwatchlogs.Client
	}

	amazonOpt func(*amazon)
)

var Conductor Amazon

func buildAmazonConductor(opts ...amazonOpt) Amazon {
	a := &amazon{}
	for _, opt := range opts {
		opt(a)
	}

	return a
}

func withConfig(cfg aws.Config) amazonOpt {
	return func(a *amazon) {
		a.amazonConfig = cfg
	}
}

func withSecretManagerClient(client *secretsmanager.Client) amazonOpt {
	return func(a *amazon) {
		a.secretManagerClient = client
	}
}

func withLambdaClient(client *lambda.Client) amazonOpt {
	return func(a *amazon) {
		a.lambdaClient = client
	}
}

func withSchedulerClient(client *scheduler.Client) amazonOpt {
	return func(a *amazon) {
		a.schedulerClient = client
	}
}

func withCloudWatchLogsClient(client *cloudwatchlogs.Client) amazonOpt {
	return func(a *amazon) {
		a.cloudWatchLogsClient = client
	}
}
