package amazonx

import (
	"testing"

	"github.com/57blocks/auto-action/server/internal/testdata"
	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestCreateSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSecretManagerClient := testdata.NewMockSecretManagerClient(ctrl)

	expectedOutput := &secretsmanager.CreateSecretOutput{
		ARN: aws.String("arn:aws:secretsmanager:us-east-1:123456789012:secret:test"),
	}
	mockSecretManagerClient.EXPECT().CreateSecret(gomock.Any(), gomock.Any()).
		Return(expectedOutput, nil)

	amazon := &amazon{
		secretManagerClient: mockSecretManagerClient,
	}
	ctx := new(gin.Context)
	output, err := amazon.CreateSecret(ctx, &secretsmanager.CreateSecretInput{})
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)
}

func TestPutResourcePolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSecretManagerClient := testdata.NewMockSecretManagerClient(ctrl)

	expectedOutput := &secretsmanager.PutResourcePolicyOutput{}
	mockSecretManagerClient.EXPECT().PutResourcePolicy(gomock.Any(), gomock.Any()).
		Return(expectedOutput, nil)

	amazon := &amazon{
		secretManagerClient: mockSecretManagerClient,
	}
	ctx := new(gin.Context)
	output, err := amazon.PutResourcePolicy(ctx, &secretsmanager.PutResourcePolicyInput{})
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)
}

func TestGetSecretValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSecretManagerClient := testdata.NewMockSecretManagerClient(ctrl)

	expectedOutput := &secretsmanager.GetSecretValueOutput{}
	mockSecretManagerClient.EXPECT().GetSecretValue(gomock.Any(), gomock.Any()).
		Return(expectedOutput, nil)

	amazon := &amazon{
		secretManagerClient: mockSecretManagerClient,
	}
	ctx := new(gin.Context)
	output, err := amazon.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{})
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)
}

func TestRemoveLambda(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLambdaClient := testdata.NewMockLambdaClient(ctrl)

	expectedOutput := &lambda.DeleteFunctionOutput{}
	mockLambdaClient.EXPECT().DeleteFunction(gomock.Any(), gomock.Any()).
		Return(expectedOutput, nil)

	amazon := &amazon{
		lambdaClient: mockLambdaClient,
	}
	ctx := new(gin.Context)
	output, err := amazon.RemoveLambda(ctx, &lambda.DeleteFunctionInput{})
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)
}

func TestRegisterLambda(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLambdaClient := testdata.NewMockLambdaClient(ctrl)

	expectedOutput := &lambda.CreateFunctionOutput{}
	mockLambdaClient.EXPECT().CreateFunction(gomock.Any(), gomock.Any()).
		Return(expectedOutput, nil)

	amazon := &amazon{
		lambdaClient: mockLambdaClient,
	}
	ctx := new(gin.Context)
	output, err := amazon.RegisterLambda(ctx, &lambda.CreateFunctionInput{})
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)
}

func TestInvokeLambda(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLambdaClient := testdata.NewMockLambdaClient(ctrl)

	expectedOutput := &lambda.InvokeOutput{}
	mockLambdaClient.EXPECT().Invoke(gomock.Any(), gomock.Any()).
		Return(expectedOutput, nil)

	amazon := &amazon{
		lambdaClient: mockLambdaClient,
	}
	ctx := new(gin.Context)
	output, err := amazon.InvokeLambda(ctx, &lambda.InvokeInput{})
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)
}

func TestBoundScheduler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSchedulerClient := testdata.NewMockSchedulerClient(ctrl)

	expectedOutput := &scheduler.CreateScheduleOutput{}
	mockSchedulerClient.EXPECT().CreateSchedule(gomock.Any(), gomock.Any()).
		Return(expectedOutput, nil)

	amazon := &amazon{
		schedulerClient: mockSchedulerClient,
	}
	ctx := new(gin.Context)
	output, err := amazon.BoundScheduler(ctx, &scheduler.CreateScheduleInput{})
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)
}

func TestGetScheduler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSchedulerClient := testdata.NewMockSchedulerClient(ctrl)

	expectedOutput := &scheduler.GetScheduleOutput{}
	mockSchedulerClient.EXPECT().GetSchedule(gomock.Any(), gomock.Any()).
		Return(expectedOutput, nil)

	amazon := &amazon{
		schedulerClient: mockSchedulerClient,
	}
	ctx := new(gin.Context)
	output, err := amazon.GetScheduler(ctx, &scheduler.GetScheduleInput{})
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)
}

func TestRemoveScheduler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSchedulerClient := testdata.NewMockSchedulerClient(ctrl)

	expectedOutput := &scheduler.DeleteScheduleOutput{}
	mockSchedulerClient.EXPECT().DeleteSchedule(gomock.Any(), gomock.Any()).
		Return(expectedOutput, nil)

	amazon := &amazon{
		schedulerClient: mockSchedulerClient,
	}
	ctx := new(gin.Context)
	output, err := amazon.RemoveScheduler(ctx, &scheduler.DeleteScheduleInput{})
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)
}

func TestDescribeLogStreams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCloudWatchLogsClient := testdata.NewMockCloudWatchLogsClient(ctrl)

	expectedOutput := &cloudwatchlogs.DescribeLogStreamsOutput{}
	mockCloudWatchLogsClient.EXPECT().DescribeLogStreams(gomock.Any(), gomock.Any()).
		Return(expectedOutput, nil)

	amazon := &amazon{
		cloudWatchLogsClient: mockCloudWatchLogsClient,
	}
	ctx := new(gin.Context)
	output, err := amazon.DescribeLogStreams(ctx, &cloudwatchlogs.DescribeLogStreamsInput{})
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)
}

func TestGetLogEvents(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCloudWatchLogsClient := testdata.NewMockCloudWatchLogsClient(ctrl)

	expectedOutput := &cloudwatchlogs.GetLogEventsOutput{}
	mockCloudWatchLogsClient.EXPECT().GetLogEvents(gomock.Any(), gomock.Any()).
		Return(expectedOutput, nil)

	amazon := &amazon{
		cloudWatchLogsClient: mockCloudWatchLogsClient,
	}
	ctx := new(gin.Context)
	output, err := amazon.GetLogEvents(ctx, &cloudwatchlogs.GetLogEventsInput{})
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)
}

func TestGetRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIamClient := testdata.NewMockIamClient(ctrl)

	expectedOutput := &iam.GetRoleOutput{}
	mockIamClient.EXPECT().GetRole(gomock.Any(), gomock.Any()).
		Return(expectedOutput, nil)

	amazon := &amazon{
		iamClient: mockIamClient,
	}
	ctx := new(gin.Context)
	output, err := amazon.GetRole(ctx, &iam.GetRoleInput{})
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)
}

func TestCreateRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIamClient := testdata.NewMockIamClient(ctrl)

	expectedOutput := &iam.CreateRoleOutput{}
	mockIamClient.EXPECT().CreateRole(gomock.Any(), gomock.Any()).
		Return(expectedOutput, nil)

	amazon := &amazon{
		iamClient: mockIamClient,
	}
	ctx := new(gin.Context)
	output, err := amazon.CreateRole(ctx, &iam.CreateRoleInput{})
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)
}

func TestPutRolePolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIamClient := testdata.NewMockIamClient(ctrl)

	expectedOutput := &iam.PutRolePolicyOutput{}
	mockIamClient.EXPECT().PutRolePolicy(gomock.Any(), gomock.Any()).
		Return(expectedOutput, nil)

	amazon := &amazon{
		iamClient: mockIamClient,
	}
	ctx := new(gin.Context)
	output, err := amazon.PutRolePolicy(ctx, &iam.PutRolePolicyInput{})
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)
}
