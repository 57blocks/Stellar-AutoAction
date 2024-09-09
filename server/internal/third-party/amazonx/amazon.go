package amazonx

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

//go:generate mockgen -destination ./amazon_mock.go -package amazonx -source amazon.go Amazon
type (
	Amazon interface {
		RegisterLambda(
			c context.Context,
			input *lambda.CreateFunctionInput,
			opts ...func(*lambda.Options),
		) (*lambda.CreateFunctionOutput, error)
		BoundScheduler(
			c context.Context,
			input *scheduler.CreateScheduleInput,
			opts ...func(*scheduler.Options),
		) (*scheduler.CreateScheduleOutput, error)
		InvokeLambda(
			c context.Context,
			input *lambda.InvokeInput,
		) (*lambda.InvokeOutput, error)
		DescribeLogStreams(
			c context.Context,
			input *cloudwatchlogs.DescribeLogStreamsInput,
		) (*cloudwatchlogs.DescribeLogStreamsOutput, error)
		GetLogEvents(
			c context.Context,
			input *cloudwatchlogs.GetLogEventsInput,
		) (*cloudwatchlogs.GetLogEventsOutput, error)
		GetSecretValue(
			c context.Context,
			input *secretsmanager.GetSecretValueInput,
		) (*secretsmanager.GetSecretValueOutput, error)
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

func buildAmazonConductor(opts ...amazonOpt) Amazon {
	a := new(amazon)
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

// Conductor implementation of Amazon
var Conductor Amazon

func (a *amazon) RegisterLambda(c context.Context, input *lambda.CreateFunctionInput, opts ...func(*lambda.Options)) (*lambda.CreateFunctionOutput,
	error) {
	return a.lambdaClient.CreateFunction(c, input, opts...)
}

func (a *amazon) BoundScheduler(c context.Context, input *scheduler.CreateScheduleInput, opts ...func(*scheduler.Options)) (*scheduler.CreateScheduleOutput, error) {
	return a.schedulerClient.CreateSchedule(c, input, opts...)
}

func (a *amazon) InvokeLambda(c context.Context, input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {
	return a.lambdaClient.Invoke(c, input)
}

func (a *amazon) DescribeLogStreams(
	c context.Context,
	input *cloudwatchlogs.DescribeLogStreamsInput,
) (*cloudwatchlogs.DescribeLogStreamsOutput, error) {
	return a.cloudWatchLogsClient.DescribeLogStreams(c, input)
}

func (a *amazon) GetLogEvents(
	c context.Context,
	input *cloudwatchlogs.GetLogEventsInput,
) (*cloudwatchlogs.GetLogEventsOutput, error) {
	return a.cloudWatchLogsClient.GetLogEvents(c, input)
}

func (a *amazon) GetSecretValue(
	c context.Context,
	input *secretsmanager.GetSecretValueInput,
) (*secretsmanager.GetSecretValueOutput, error) {
	return a.secretManagerClient.GetSecretValue(c, input)
}
