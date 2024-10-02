package amazonx

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

//go:generate mockgen -destination ../../testdata/amazon_mock.go -package testdata -source amazon.go Amazon
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
		GetScheduler(
			c context.Context,
			input *scheduler.GetScheduleInput,
		) (*scheduler.GetScheduleOutput, error)
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
		RemoveLambda(
			c context.Context,
			input *lambda.DeleteFunctionInput,
		) (*lambda.DeleteFunctionOutput, error)
		RemoveScheduler(
			c context.Context,
			input *scheduler.DeleteScheduleInput,
		) (*scheduler.DeleteScheduleOutput, error)
		GetRole(
			c context.Context,
			input *iam.GetRoleInput,
		) (*iam.GetRoleOutput, error)
		CreateRole(
			c context.Context,
			input *iam.CreateRoleInput,
		) (*iam.CreateRoleOutput, error)
		PutRolePolicy(
			c context.Context,
			input *iam.PutRolePolicyInput,
		) (*iam.PutRolePolicyOutput, error)
		CreateSecret(
			c context.Context,
			input *secretsmanager.CreateSecretInput,
		) (*secretsmanager.CreateSecretOutput, error)
		PutResourcePolicy(
			c context.Context,
			input *secretsmanager.PutResourcePolicyInput,
		) (*secretsmanager.PutResourcePolicyOutput, error)
	}

	SecretManagerClient interface {
		GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
		CreateSecret(ctx context.Context, params *secretsmanager.CreateSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.CreateSecretOutput, error)
		PutResourcePolicy(ctx context.Context, params *secretsmanager.PutResourcePolicyInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.PutResourcePolicyOutput, error)
	}

	LambdaClient interface {
		DeleteFunction(ctx context.Context, params *lambda.DeleteFunctionInput, optFns ...func(*lambda.Options)) (*lambda.DeleteFunctionOutput, error)
		Invoke(ctx context.Context, params *lambda.InvokeInput, optFns ...func(*lambda.Options)) (*lambda.InvokeOutput, error)
		CreateFunction(ctx context.Context, params *lambda.CreateFunctionInput, optFns ...func(*lambda.Options)) (*lambda.CreateFunctionOutput, error)
	}

	SchedulerClient interface {
		CreateSchedule(ctx context.Context, params *scheduler.CreateScheduleInput, optFns ...func(*scheduler.Options)) (*scheduler.CreateScheduleOutput, error)
		DeleteSchedule(ctx context.Context, params *scheduler.DeleteScheduleInput, optFns ...func(*scheduler.Options)) (*scheduler.DeleteScheduleOutput, error)
		GetSchedule(ctx context.Context, params *scheduler.GetScheduleInput, optFns ...func(*scheduler.Options)) (*scheduler.GetScheduleOutput, error)
	}

	CloudWatchLogsClient interface {
		DescribeLogStreams(ctx context.Context, params *cloudwatchlogs.DescribeLogStreamsInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogStreamsOutput, error)
		GetLogEvents(ctx context.Context, params *cloudwatchlogs.GetLogEventsInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.GetLogEventsOutput, error)
	}

	IamClient interface {
		GetRole(ctx context.Context, params *iam.GetRoleInput, optFns ...func(*iam.Options)) (*iam.GetRoleOutput, error)
		CreateRole(ctx context.Context, params *iam.CreateRoleInput, optFns ...func(*iam.Options)) (*iam.CreateRoleOutput, error)
		PutRolePolicy(ctx context.Context, params *iam.PutRolePolicyInput, optFns ...func(*iam.Options)) (*iam.PutRolePolicyOutput, error)
	}

	amazon struct {
		amazonConfig         aws.Config
		secretManagerClient  SecretManagerClient
		lambdaClient         LambdaClient
		schedulerClient      SchedulerClient
		cloudWatchLogsClient CloudWatchLogsClient
		iamClient            IamClient
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

func withSecretManagerClient(client SecretManagerClient) amazonOpt {
	return func(a *amazon) {
		a.secretManagerClient = client
	}
}

func withLambdaClient(client LambdaClient) amazonOpt {
	return func(a *amazon) {
		a.lambdaClient = client
	}
}

func withSchedulerClient(client SchedulerClient) amazonOpt {
	return func(a *amazon) {
		a.schedulerClient = client
	}
}

func withCloudWatchLogsClient(client CloudWatchLogsClient) amazonOpt {
	return func(a *amazon) {
		a.cloudWatchLogsClient = client
	}
}

func withIamClient(client IamClient) amazonOpt {
	return func(a *amazon) {
		a.iamClient = client
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

func (a *amazon) GetScheduler(c context.Context, input *scheduler.GetScheduleInput) (*scheduler.GetScheduleOutput, error) {
	return a.schedulerClient.GetSchedule(c, input)
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

func (a *amazon) RemoveLambda(
	c context.Context,
	input *lambda.DeleteFunctionInput,
) (*lambda.DeleteFunctionOutput, error) {
	return a.lambdaClient.DeleteFunction(c, input)
}

func (a *amazon) RemoveScheduler(
	c context.Context,
	input *scheduler.DeleteScheduleInput,
) (*scheduler.DeleteScheduleOutput, error) {
	return a.schedulerClient.DeleteSchedule(c, input)
}

func (a *amazon) GetRole(
	c context.Context,
	input *iam.GetRoleInput,
) (*iam.GetRoleOutput, error) {
	return a.iamClient.GetRole(c, input)
}

func (a *amazon) CreateRole(
	c context.Context,
	input *iam.CreateRoleInput,
) (*iam.CreateRoleOutput, error) {
	return a.iamClient.CreateRole(c, input)
}

func (a *amazon) PutRolePolicy(
	c context.Context,
	input *iam.PutRolePolicyInput,
) (*iam.PutRolePolicyOutput, error) {
	return a.iamClient.PutRolePolicy(c, input)
}

func (a *amazon) CreateSecret(
	c context.Context,
	input *secretsmanager.CreateSecretInput,
) (*secretsmanager.CreateSecretOutput, error) {
	return a.secretManagerClient.CreateSecret(c, input)
}

func (a *amazon) PutResourcePolicy(
	c context.Context,
	input *secretsmanager.PutResourcePolicyInput,
) (*secretsmanager.PutResourcePolicyOutput, error) {
	return a.secretManagerClient.PutResourcePolicy(c, input)
}
