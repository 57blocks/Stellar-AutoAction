package lambda

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/constant"
	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/57blocks/auto-action/server/internal/testdata"
	"github.com/57blocks/auto-action/server/internal/third-party/logx"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	lambTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	scheTypes "github.com/aws/aws-sdk-go-v2/service/scheduler/types"
	"github.com/aws/smithy-go/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

// Before test, setup log and config
func TestMain(m *testing.M) {
	os.Setenv("LAMBDA_MAX", "2")
	config.Setup("../../config/")
	testConfig := config.Configuration{
		Log: config.Log{
			Level:    "debug",
			Encoding: "json",
		},
	}
	logx.Setup(&testConfig)

	os.Exit(m.Run())
}

func TestInfoSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLambRepo := testdata.NewMockLambda(ctrl)
	mockOAuthRepo := testdata.NewMockOAuth(ctrl)

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimSub.Str(), "account_name")
	request := &dto.ReqURILambda{Lambda: "name/arn"}
	distinguisher := "name/arn"
	accountID := uint64(123)

	now := time.Now().UTC()
	sha256 := "E2X2ZXxocZcefGFb8lu2QnbYV8higgV2yYcJSwPLAA4="
	schARN := "arn:aws:scheduler:...:schedule/default/transfer"

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, "account_name").Times(1).
		Return(&dto.RespUser{ID: accountID}, nil)

	mockLambRepo.EXPECT().LambdaInfo(ctx, accountID, distinguisher).Times(1).
		Return(
			&dto.RespInfo{
				CodeSHA256:   sha256,
				Description:  "",
				FunctionArn:  "arn:aws:lambda:us-east-2:123340007534:function:epic-v3n-transfer",
				FunctionName: "epic-v3n-transfer",
				Handler:      "transfer.handler",
				RevisionID:   "3c9a3513-5e43-419e-ae5b-aeeb459e44e3",
				Role:         "arn:aws:iam::123456789012:role/LambdaExecutionRole",
				Runtime:      "nodejs20.x",
				Version:      "$LATEST",
				Scheduler: dto.Scheduler{
					Expression:  "rate(1 minutes)",
					ScheduleArn: schARN,
				},
				CreatedAt: &now,
				UpdatedAt: &now,
			}, nil,
		)

	cd := &service{
		lambdaRepo: mockLambRepo,
		oauthRepo:  mockOAuthRepo,
	}
	info, err := cd.Info(ctx, request)
	assert.NoError(t, err)
	assert.Equal(t, sha256, info.CodeSHA256)
	assert.Equal(t, schARN, info.Scheduler.ScheduleArn)
}

func TestInfoUserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOAuthRepo := testdata.NewMockOAuth(ctrl)

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimSub.Str(), "account_name")
	request := &dto.ReqURILambda{Lambda: "name/arn"}

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, "account_name").Times(1).
		Return(nil, errorx.NotFound("user not found"))

	cd := &service{
		oauthRepo: mockOAuthRepo,
	}
	info, err := cd.Info(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, errorx.NotFound("user not found"), err)
	assert.Nil(t, info)
}

func TestInfoLambdaNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLambRepo := testdata.NewMockLambda(ctrl)
	mockOAuthRepo := testdata.NewMockOAuth(ctrl)

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimSub.Str(), "account_name")
	request := &dto.ReqURILambda{Lambda: "name/arn"}
	accountID := uint64(123)

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, "account_name").Times(1).
		Return(&dto.RespUser{ID: accountID}, nil)

	mockLambRepo.EXPECT().LambdaInfo(ctx, accountID, "name/arn").Times(1).
		Return(nil, errorx.NotFound("lambda not found"))

	cd := &service{
		lambdaRepo: mockLambRepo,
		oauthRepo:  mockOAuthRepo,
	}

	info, err := cd.Info(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, errorx.NotFound("lambda not found"), err)
	assert.Nil(t, info)
}

func TestInvokeSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLambRepo := testdata.NewMockLambda(ctrl)
	mockAmazon := testdata.NewMockAmazon(ctrl)
	mockOAuthRepo := testdata.NewMockOAuth(ctrl)

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimSub.Str(), "account_name")
	ctx.Set(constant.ClaimIss.Str(), "org_name")
	request := &dto.ReqInvoke{
		Lambda:  "name/arn",
		Payload: "{\"foo\":\"bar\"}",
	}
	accountID := uint64(123)

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, "account_name").Times(1).
		Return(&dto.RespUser{ID: accountID}, nil)

	mockLambRepo.EXPECT().LambdaInfo(ctx, accountID, "name/arn").Times(1).
		Return(&dto.RespInfo{
			FunctionName: "name/arn",
		}, nil)

	encodedLogResult := base64.StdEncoding.EncodeToString([]byte("foo"))
	encodedPayload := base64.StdEncoding.EncodeToString([]byte("bar"))
	awsInputPayload, _ := json.Marshal(map[string]interface{}{
		"organization": "org_name",
		"account":      "account_name",
		"foo":          "bar",
	})
	mockAmazon.EXPECT().InvokeLambda(ctx, &lambda.InvokeInput{
		FunctionName: aws.String("name/arn"),
		LogType:      lambTypes.LogTypeTail,
		Payload:      awsInputPayload,
	}).Times(1).
		Return(&lambda.InvokeOutput{
			LogResult: aws.String(encodedLogResult),
			Payload:   []byte(encodedPayload),
		}, nil)

	cd := &service{
		lambdaRepo: mockLambRepo,
		amazon:     mockAmazon,
		oauthRepo:  mockOAuthRepo,
	}

	invoke, err := cd.Invoke(ctx, request)
	assert.NoError(t, err)
	assert.Equal(t, "foo", *invoke.LogResult)
	assert.Equal(t, []byte("bar"), invoke.Payload)
}

func TestInvokeUserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOAuthRepo := testdata.NewMockOAuth(ctrl)

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimSub.Str(), "account_name")
	request := &dto.ReqInvoke{
		Lambda: "name/arn",
	}

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, "account_name").Times(1).
		Return(nil, errorx.NotFound("user not found"))

	cd := &service{
		oauthRepo: mockOAuthRepo,
	}

	invoke, err := cd.Invoke(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, errorx.NotFound("user not found"), err)
	assert.Nil(t, invoke)
}

func TestInvokeLambdaNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLambRepo := testdata.NewMockLambda(ctrl)
	mockOAuthRepo := testdata.NewMockOAuth(ctrl)

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimSub.Str(), "account_name")
	request := &dto.ReqInvoke{
		Lambda: "name/arn",
	}
	accountID := uint64(123)

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, "account_name").Times(1).
		Return(&dto.RespUser{ID: accountID}, nil)

	mockLambRepo.EXPECT().LambdaInfo(ctx, accountID, "name/arn").Times(1).
		Return(nil, errorx.NotFound("lambda not found"))

	cd := &service{
		lambdaRepo: mockLambRepo,
		oauthRepo:  mockOAuthRepo,
	}

	invoke, err := cd.Invoke(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, errorx.NotFound("lambda not found"), err)
	assert.Nil(t, invoke)
}

func TestInvokePayloadInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLambRepo := testdata.NewMockLambda(ctrl)
	mockOAuthRepo := testdata.NewMockOAuth(ctrl)
	ctx := new(gin.Context)
	ctx.Set(constant.ClaimSub.Str(), "account_name")
	ctx.Set(constant.ClaimIss.Str(), "org_name")
	request := &dto.ReqInvoke{
		Lambda:  "name/arn",
		Payload: "error payload",
	}
	accountID := uint64(123)

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, "account_name").Times(1).
		Return(&dto.RespUser{ID: accountID}, nil)

	mockLambRepo.EXPECT().LambdaInfo(ctx, accountID, "name/arn").Times(1).
		Return(&dto.RespInfo{
			FunctionName: "name/arn",
		}, nil)

	cd := &service{
		lambdaRepo: mockLambRepo,
		oauthRepo:  mockOAuthRepo,
	}

	invoke, err := cd.Invoke(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, errorx.Internal("failed to unmarshal payload: invalid character 'e' looking for beginning of value"), err)
	assert.Nil(t, invoke)
}

func TestInvokeLambdaError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLambRepo := testdata.NewMockLambda(ctrl)
	mockAmazon := testdata.NewMockAmazon(ctrl)
	mockOAuthRepo := testdata.NewMockOAuth(ctrl)

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimSub.Str(), "account_name")
	ctx.Set(constant.ClaimIss.Str(), "org_name")
	request := &dto.ReqInvoke{
		Lambda:  "name/arn",
		Payload: "{\"foo\":\"bar\"}",
	}
	accountID := uint64(123)

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, "account_name").Times(1).
		Return(&dto.RespUser{ID: accountID}, nil)

	mockLambRepo.EXPECT().LambdaInfo(ctx, accountID, "name/arn").Times(1).
		Return(&dto.RespInfo{
			FunctionName: "name/arn",
		}, nil)

	awsInputPayload, _ := json.Marshal(map[string]interface{}{
		"organization": "org_name",
		"account":      "account_name",
		"foo":          "bar",
	})
	mockAmazon.EXPECT().InvokeLambda(ctx, &lambda.InvokeInput{
		FunctionName: aws.String("name/arn"),
		LogType:      lambTypes.LogTypeTail,
		Payload:      awsInputPayload,
	}).Times(1).
		Return(nil, errorx.Internal("failed to invoke lambda"))

	cd := &service{
		lambdaRepo: mockLambRepo,
		amazon:     mockAmazon,
		oauthRepo:  mockOAuthRepo,
	}

	invoke, err := cd.Invoke(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, errorx.Internal("failed to invoke lambda: name/arn, error: failed to invoke lambda"), err)
	assert.Nil(t, invoke)
}

func TestInvokeDecodeLogResultError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLambRepo := testdata.NewMockLambda(ctrl)
	mockAmazon := testdata.NewMockAmazon(ctrl)
	mockOAuthRepo := testdata.NewMockOAuth(ctrl)

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimSub.Str(), "account_name")
	ctx.Set(constant.ClaimIss.Str(), "org_name")
	request := &dto.ReqInvoke{
		Lambda:  "name/arn",
		Payload: "{\"foo\":\"bar\"}",
	}
	accountID := uint64(123)

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, "account_name").Times(1).
		Return(&dto.RespUser{ID: accountID}, nil)

	mockLambRepo.EXPECT().LambdaInfo(ctx, accountID, "name/arn").Times(1).
		Return(&dto.RespInfo{
			FunctionName: "name/arn",
		}, nil)

	awsInputPayload, _ := json.Marshal(map[string]interface{}{
		"organization": "org_name",
		"account":      "account_name",
		"foo":          "bar",
	})
	mockAmazon.EXPECT().InvokeLambda(ctx, &lambda.InvokeInput{
		FunctionName: aws.String("name/arn"),
		LogType:      lambTypes.LogTypeTail,
		Payload:      awsInputPayload,
	}).Times(1).
		Return(&lambda.InvokeOutput{
			LogResult: aws.String("hello"),
		}, nil)

	cd := &service{
		lambdaRepo: mockLambRepo,
		amazon:     mockAmazon,
		oauthRepo:  mockOAuthRepo,
	}

	invoke, err := cd.Invoke(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, errorx.Internal("failed to decode log result: illegal base64 data at input byte 4"), err)
	assert.Nil(t, invoke)
}

func TestInvokeDecodePayloadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLambRepo := testdata.NewMockLambda(ctrl)
	mockAmazon := testdata.NewMockAmazon(ctrl)
	mockOAuthRepo := testdata.NewMockOAuth(ctrl)
	ctx := new(gin.Context)
	ctx.Set(constant.ClaimSub.Str(), "account_name")
	ctx.Set(constant.ClaimIss.Str(), "org_name")
	request := &dto.ReqInvoke{
		Lambda:  "name/arn",
		Payload: "{\"foo\":\"bar\"}",
	}
	accountID := uint64(123)

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, "account_name").Times(1).
		Return(&dto.RespUser{ID: accountID}, nil)

	mockLambRepo.EXPECT().LambdaInfo(ctx, accountID, "name/arn").Times(1).
		Return(&dto.RespInfo{
			FunctionName: "name/arn",
		}, nil)

	awsInputPayload, _ := json.Marshal(map[string]interface{}{
		"organization": "org_name",
		"account":      "account_name",
		"foo":          "bar",
	})
	encodedLogResult := base64.StdEncoding.EncodeToString([]byte("foo"))
	mockAmazon.EXPECT().InvokeLambda(ctx, &lambda.InvokeInput{
		FunctionName: aws.String("name/arn"),
		LogType:      lambTypes.LogTypeTail,
		Payload:      awsInputPayload,
	}).Times(1).
		Return(&lambda.InvokeOutput{
			LogResult: aws.String(encodedLogResult),
			Payload:   []byte("hello"),
		}, nil)

	cd := &service{
		lambdaRepo: mockLambRepo,
		amazon:     mockAmazon,
		oauthRepo:  mockOAuthRepo,
	}

	invoke, err := cd.Invoke(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, errorx.Internal("failed to decode payload: illegal base64 data at input byte 4"), err)
	assert.Nil(t, invoke)
}

func TestListSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLambRepo := testdata.NewMockLambda(ctrl)
	mockOAuthRepo := testdata.NewMockOAuth(ctrl)

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimSub.Str(), "account_name")
	accountID := uint64(123)
	now := time.Now().UTC()

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, "account_name").Times(1).
		Return(&dto.RespUser{ID: accountID}, nil)

	expectedLamb := dto.RespInfo{
		FunctionName: "name/arn",
		FunctionArn:  "arn:aws:lambda:us-east-2:123340007534:function:epic-v3n-transfer",
		Description:  "",
		CreatedAt:    &now,
	}
	mockLambRepo.EXPECT().FindByAccount(ctx, accountID).Times(1).
		Return([]*dto.RespInfo{&expectedLamb}, nil)

	cd := &service{
		lambdaRepo: mockLambRepo,
		oauthRepo:  mockOAuthRepo,
	}

	expectedRespInList := &dto.RespInList{
		FunctionName: expectedLamb.FunctionName,
		FunctionArn:  expectedLamb.FunctionArn,
		Description:  expectedLamb.Description,
		CreatedAt:    expectedLamb.CreatedAt,
	}
	list, err := cd.List(ctx, false)
	assert.NoError(t, err)
	assert.Equal(t, expectedRespInList, list.([]*dto.RespInList)[0])
}

func TestListFullSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLambRepo := testdata.NewMockLambda(ctrl)
	mockOAuthRepo := testdata.NewMockOAuth(ctrl)

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimSub.Str(), "account_name")
	accountID := uint64(123)
	now := time.Now().UTC()

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, "account_name").Times(1).
		Return(&dto.RespUser{ID: accountID}, nil)

	expectedLamb := dto.RespInfo{
		FunctionName: "name/arn",
		FunctionArn:  "arn:aws:lambda:us-east-2:123340007534:function:epic-v3n-transfer",
		Description:  "",
		CreatedAt:    &now,
	}
	mockLambRepo.EXPECT().FindByAccount(ctx, accountID).Times(1).
		Return([]*dto.RespInfo{&expectedLamb}, nil)

	cd := &service{
		lambdaRepo: mockLambRepo,
		oauthRepo:  mockOAuthRepo,
	}

	list, err := cd.List(ctx, true)
	assert.NoError(t, err)
	assert.Equal(t, &expectedLamb, list.([]*dto.RespInfo)[0])
}

func TestListUserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOAuthRepo := testdata.NewMockOAuth(ctrl)

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimSub.Str(), "account_name")

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, "account_name").Times(1).
		Return(nil, errorx.NotFound("user not found"))

	cd := &service{
		oauthRepo: mockOAuthRepo,
	}

	list, err := cd.List(ctx, false)
	assert.Error(t, err)
	assert.Equal(t, errorx.NotFound("user not found"), err)
	assert.Nil(t, list)
}

func TestListLambdaError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLambRepo := testdata.NewMockLambda(ctrl)
	mockOAuthRepo := testdata.NewMockOAuth(ctrl)

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimSub.Str(), "account_name")
	accountID := uint64(123)

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, "account_name").Times(1).
		Return(&dto.RespUser{ID: accountID}, nil)

	mockLambRepo.EXPECT().FindByAccount(ctx, accountID).Times(1).
		Return(nil, errorx.Internal("failed to find lambda"))

	cd := &service{
		lambdaRepo: mockLambRepo,
		oauthRepo:  mockOAuthRepo,
	}

	list, err := cd.List(ctx, false)
	assert.Error(t, err)
	assert.Equal(t, errorx.Internal("failed to find lambda"), err)
	assert.Nil(t, list)
}

func TestRegisterSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLambRepo := testdata.NewMockLambda(ctrl)
	mockOAuthRepo := testdata.NewMockOAuth(ctrl)
	mockAmazon := testdata.NewMockAmazon(ctrl)

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimSub.Str(), "account_name")
	ctx.Set(constant.ClaimIss.Str(), "org_name")
	request := &dto.ReqRegister{
		Expression: "rate(minutes)",
		Files: []*dto.ReqFile{
			{
				Name:  "file1",
				Bytes: []byte("file1"),
			},
		},
	}
	accountID := uint64(123)
	roleARN := "arn:aws:iam::123456789012:role/AA-org_name-account_name-Role"
	functionARN := "arn:aws:lambda:us-east-2:123456789012:function:file1"
	scheduleARN := "arn:aws:scheduler:us-east-2:123456789012:schedule/default/file1"

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, &dto.ReqOrgAcn{
		OrgName: "org_name",
		AcnName: "account_name",
	}).Times(1).
		Return(&dto.RespUser{ID: accountID}, nil)

	expectedLambs := make([]*dto.RespInfo, 0)
	mockLambRepo.EXPECT().FindByAccount(ctx, accountID).Times(1).
		Return(expectedLambs, nil)

	mockAmazon.EXPECT().GetRole(ctx, &iam.GetRoleInput{
		RoleName: aws.String("AA-org_name-account_name-Role"),
	}).Times(1).
		Return(&iam.GetRoleOutput{
			Role: &iamTypes.Role{
				Arn: aws.String(roleARN),
			},
		}, nil)

	mockAmazon.EXPECT().RegisterLambda(ctx, gomock.Any(), gomock.Any()).Times(1).
		DoAndReturn(func(ctx *gin.Context, input *lambda.CreateFunctionInput, _ ...func(*lambda.Options)) (*lambda.CreateFunctionOutput, error) {
			assert.Equal(t, []byte("file1"), input.Code.ZipFile)
			assert.Equal(t, aws.String("org_name-account_name-file1"), input.FunctionName)
			assert.Equal(t, map[string]string{"ENV_AWS_REGION": "us-east-2"}, input.Environment.Variables)
			assert.Equal(t, aws.String(roleARN), input.Role)
			assert.Equal(t, lambTypes.RuntimeNodejs20x, input.Runtime)
			assert.Equal(t, aws.Int32(30), input.Timeout)
			assert.Equal(t, aws.String("file1.handler"), input.Handler)
			assert.Equal(t, lambTypes.PackageTypeZip, input.PackageType)
			assert.Equal(t, false, input.Publish)

			return &lambda.CreateFunctionOutput{
				FunctionName: aws.String("org_name-account_name-file1"),
				FunctionArn:  aws.String(functionARN),
				Runtime:      lambTypes.RuntimeNodejs20x,
				Handler:      aws.String("file1.handler"),
				Version:      aws.String("$LATEST"),
				Timeout:      aws.Int32(30),
				Role:         aws.String(roleARN),
				Description:  aws.String(""),
				CodeSha256:   aws.String("E2X2ZXxocZcefGFb8lu2QnbYV8higgV2yYcJSwPLAA4="),
				RevisionId:   aws.String("3c9a3513-5e43-419e-ae5b-aeeb459e44e3"),
			}, nil
		})

	mockAmazon.EXPECT().BoundScheduler(ctx, gomock.Any()).Times(1).
		DoAndReturn(func(_ *gin.Context, input *scheduler.CreateScheduleInput, _ ...func(*scheduler.Options)) (*scheduler.CreateScheduleOutput, error) {
			assert.Equal(t, scheTypes.FlexibleTimeWindowModeOff, input.FlexibleTimeWindow.Mode)
			assert.Equal(t, aws.String("org_name-account_name-file1"), input.Name)
			assert.Equal(t, aws.String("rate(minutes)"), input.ScheduleExpression)
			assert.Equal(t, aws.String(functionARN), input.Target.Arn)

			return &scheduler.CreateScheduleOutput{
				ScheduleArn: aws.String(scheduleARN),
			}, nil
		})

	mockLambRepo.EXPECT().PersistRegResult(ctx, gomock.Any()).Times(1).
		Return(nil)

	cd := &service{
		lambdaRepo: mockLambRepo,
		oauthRepo:  mockOAuthRepo,
		amazon:     mockAmazon,
	}

	expectedResp := []*dto.RespRegister{
		{
			Lambda: &dto.RespLamBrief{
				Name:    "org_name-account_name-file1",
				Arn:     functionARN,
				Runtime: "nodejs20.x",
				Handler: "file1.handler",
				Version: "$LATEST",
			},
			Scheduler: &dto.RespSchBrief{
				Arn:            scheduleARN,
				BoundLambdaArn: functionARN,
				Name:           "org_name-account_name-file1",
			},
		},
	}
	register, err := cd.Register(ctx, request)
	assert.NoError(t, err)
	assert.Equal(t, expectedResp, register)
}

func TestRegisterUserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOAuthRepo := testdata.NewMockOAuth(ctrl)

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimSub.Str(), "account_name")
	ctx.Set(constant.ClaimIss.Str(), "org_name")
	request := &dto.ReqRegister{
		Expression: "rate(minutes)",
		Files: []*dto.ReqFile{
			{
				Name:  "file1",
				Bytes: []byte("file1"),
			},
		},
	}

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, &dto.ReqOrgAcn{
		OrgName: "org_name",
		AcnName: "account_name",
	}).Times(1).
		Return(nil, errorx.NotFound("user not found"))

	cd := &service{
		oauthRepo: mockOAuthRepo,
	}

	register, err := cd.Register(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, errorx.NotFound("user not found"), err)
	assert.Nil(t, register)
}

func TestRegisterMaxLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOAuthRepo := testdata.NewMockOAuth(ctrl)
	mockLambRepo := testdata.NewMockLambda(ctrl)
	ctx := new(gin.Context)
	ctx.Set(constant.ClaimSub.Str(), "account_name")
	ctx.Set(constant.ClaimIss.Str(), "org_name")
	accountID := uint64(123)
	request := &dto.ReqRegister{
		Expression: "rate(minutes)",
		Files: []*dto.ReqFile{
			{
				Name:  "file1",
				Bytes: []byte("file1"),
			},
			{
				Name:  "file2",
				Bytes: []byte("file2"),
			},
			{
				Name:  "file3",
				Bytes: []byte("file3"),
			},
		},
	}

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, &dto.ReqOrgAcn{
		OrgName: "org_name",
		AcnName: "account_name",
	}).Times(1).
		Return(&dto.RespUser{ID: accountID}, nil)

	expectedLambs := make([]*dto.RespInfo, 0)
	mockLambRepo.EXPECT().FindByAccount(ctx, accountID).Times(1).
		Return(expectedLambs, nil)

	cd := &service{
		lambdaRepo: mockLambRepo,
		oauthRepo:  mockOAuthRepo,
	}

	register, err := cd.Register(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, errorx.BadRequest("the number of lambdas is limited to 2"), err)
	assert.Nil(t, register)
}

func TestRegisterGetRoleError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOAuthRepo := testdata.NewMockOAuth(ctrl)
	mockLambRepo := testdata.NewMockLambda(ctrl)
	mockAmazon := testdata.NewMockAmazon(ctrl)

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimSub.Str(), "account_name")
	ctx.Set(constant.ClaimIss.Str(), "org_name")
	accountID := uint64(123)
	request := &dto.ReqRegister{
		Expression: "rate(minutes)",
		Files: []*dto.ReqFile{
			{
				Name:  "file1",
				Bytes: []byte("file1"),
			},
		},
	}

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, &dto.ReqOrgAcn{
		OrgName: "org_name",
		AcnName: "account_name",
	}).Times(1).
		Return(&dto.RespUser{ID: accountID}, nil)

	expectedLambs := make([]*dto.RespInfo, 0)
	mockLambRepo.EXPECT().FindByAccount(ctx, accountID).Times(1).
		Return(expectedLambs, nil)

	mockAmazon.EXPECT().GetRole(ctx, &iam.GetRoleInput{
		RoleName: aws.String("AA-org_name-account_name-Role"),
	}).Times(1).
		Return(nil, errorx.Internal("failed to get role"))

	cd := &service{
		lambdaRepo: mockLambRepo,
		oauthRepo:  mockOAuthRepo,
		amazon:     mockAmazon,
	}

	register, err := cd.Register(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, errorx.Internal("get role arn error: failed to get role"), err)
	assert.Nil(t, register)
}

func TestRegisterRegisterLambdaError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOAuthRepo := testdata.NewMockOAuth(ctrl)
	mockLambRepo := testdata.NewMockLambda(ctrl)
	mockAmazon := testdata.NewMockAmazon(ctrl)

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimSub.Str(), "account_name")
	ctx.Set(constant.ClaimIss.Str(), "org_name")
	request := &dto.ReqRegister{
		Expression: "rate(minutes)",
		Files: []*dto.ReqFile{
			{
				Name:  "file1",
				Bytes: []byte("file1"),
			},
		},
	}
	accountID := uint64(123)
	roleARN := "arn:aws:iam::123456789012:role/AA-org_name-account_name-Role"

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, &dto.ReqOrgAcn{
		OrgName: "org_name",
		AcnName: "account_name",
	}).Times(1).
		Return(&dto.RespUser{ID: accountID}, nil)

	expectedLambs := make([]*dto.RespInfo, 0)
	mockLambRepo.EXPECT().FindByAccount(ctx, accountID).Times(1).
		Return(expectedLambs, nil)

	mockAmazon.EXPECT().GetRole(ctx, &iam.GetRoleInput{
		RoleName: aws.String("AA-org_name-account_name-Role"),
	}).Times(1).
		Return(&iam.GetRoleOutput{
			Role: &iamTypes.Role{
				Arn: aws.String(roleARN),
			},
		}, nil)

	mockAmazon.EXPECT().RegisterLambda(ctx, gomock.Any(), gomock.Any()).Times(1).
		Return(nil, errorx.Internal("failed to register lambda"))

	cd := &service{
		lambdaRepo: mockLambRepo,
		oauthRepo:  mockOAuthRepo,
		amazon:     mockAmazon,
	}

	register, err := cd.Register(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, errorx.Internal("failed to register lambda: file1, err: failed to register lambda"), err)
	assert.Nil(t, register)
}

func TestRegisterRegisterSchedulerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOAuthRepo := testdata.NewMockOAuth(ctrl)
	mockLambRepo := testdata.NewMockLambda(ctrl)
	mockAmazon := testdata.NewMockAmazon(ctrl)

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimSub.Str(), "account_name")
	ctx.Set(constant.ClaimIss.Str(), "org_name")
	request := &dto.ReqRegister{
		Expression: "rate(minutes)",
		Files: []*dto.ReqFile{
			{
				Name:  "file1",
				Bytes: []byte("file1"),
			},
		},
	}
	accountID := uint64(123)
	roleARN := "arn:aws:iam::123456789012:role/AA-org_name-account_name-Role"
	functionARN := "arn:aws:lambda:us-east-2:123456789012:function:file1"

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, &dto.ReqOrgAcn{
		OrgName: "org_name",
		AcnName: "account_name",
	}).Times(1).
		Return(&dto.RespUser{ID: accountID}, nil)

	expectedLambs := make([]*dto.RespInfo, 0)
	mockLambRepo.EXPECT().FindByAccount(ctx, accountID).Times(1).
		Return(expectedLambs, nil)

	mockAmazon.EXPECT().GetRole(ctx, &iam.GetRoleInput{
		RoleName: aws.String("AA-org_name-account_name-Role"),
	}).Times(1).
		Return(&iam.GetRoleOutput{
			Role: &iamTypes.Role{
				Arn: aws.String(roleARN),
			},
		}, nil)

	mockAmazon.EXPECT().RegisterLambda(ctx, gomock.Any(), gomock.Any()).Times(1).
		Return(&lambda.CreateFunctionOutput{
			FunctionName: aws.String("org_name-account_name-file1"),
			FunctionArn:  aws.String(functionARN),
			Runtime:      lambTypes.RuntimeNodejs20x,
			Handler:      aws.String("file1.handler"),
			Version:      aws.String("$LATEST"),
			Timeout:      aws.Int32(30),
			Role:         aws.String(roleARN),
			Description:  aws.String(""),
			CodeSha256:   aws.String("E2X2ZXxocZcefGFb8lu2QnbYV8higgV2yYcJSwPLAA4="),
			RevisionId:   aws.String("3c9a3513-5e43-419e-ae5b-aeeb459e44e3"),
		}, nil)

	mockAmazon.EXPECT().BoundScheduler(ctx, gomock.Any()).Times(1).
		Return(nil, errorx.Internal("failed to bound scheduler"))

	cd := &service{
		lambdaRepo: mockLambRepo,
		oauthRepo:  mockOAuthRepo,
		amazon:     mockAmazon,
	}

	register, err := cd.Register(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, errorx.Internal("failed to bound scheduler: org_name-account_name-file1, err: failed to bound scheduler"), err)
	assert.Nil(t, register)
}

func TestRegisterPersistRegisterResultsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOAuthRepo := testdata.NewMockOAuth(ctrl)
	mockLambRepo := testdata.NewMockLambda(ctrl)
	mockAmazon := testdata.NewMockAmazon(ctrl)

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimSub.Str(), "account_name")
	ctx.Set(constant.ClaimIss.Str(), "org_name")
	request := &dto.ReqRegister{
		Expression: "rate(minutes)",
		Files: []*dto.ReqFile{
			{
				Name:  "file1",
				Bytes: []byte("file1"),
			},
		},
	}
	accountID := uint64(123)
	roleARN := "arn:aws:iam::123456789012:role/AA-org_name-account_name-Role"
	functionARN := "arn:aws:lambda:us-east-2:123456789012:function:file1"
	scheduleARN := "arn:aws:scheduler:us-east-2:123456789012:schedule/default/file1"

	mockOAuthRepo.EXPECT().FindUserByOrgAcn(ctx, &dto.ReqOrgAcn{
		OrgName: "org_name",
		AcnName: "account_name",
	}).Times(1).
		Return(&dto.RespUser{ID: accountID}, nil)

	expectedLambs := make([]*dto.RespInfo, 0)
	mockLambRepo.EXPECT().FindByAccount(ctx, accountID).Times(1).
		Return(expectedLambs, nil)

	mockAmazon.EXPECT().GetRole(ctx, &iam.GetRoleInput{
		RoleName: aws.String("AA-org_name-account_name-Role"),
	}).Times(1).
		Return(&iam.GetRoleOutput{
			Role: &iamTypes.Role{
				Arn: aws.String(roleARN),
			},
		}, nil)

	mockAmazon.EXPECT().RegisterLambda(ctx, gomock.Any(), gomock.Any()).Times(1).
		Return(&lambda.CreateFunctionOutput{
			FunctionName: aws.String("org_name-account_name-file1"),
			FunctionArn:  aws.String(functionARN),
			Runtime:      lambTypes.RuntimeNodejs20x,
			Handler:      aws.String("file1.handler"),
			Version:      aws.String("$LATEST"),
			Timeout:      aws.Int32(30),
			Role:         aws.String(roleARN),
			Description:  aws.String(""),
			CodeSha256:   aws.String("E2X2ZXxocZcefGFb8lu2QnbYV8higgV2yYcJSwPLAA4="),
			RevisionId:   aws.String("3c9a3513-5e43-419e-ae5b-aeeb459e44e3"),
		}, nil)

	mockAmazon.EXPECT().BoundScheduler(ctx, gomock.Any()).Times(1).
		Return(&scheduler.CreateScheduleOutput{
			ScheduleArn: aws.String(scheduleARN),
		}, nil)

	mockLambRepo.EXPECT().PersistRegResult(ctx, gomock.Any()).Times(1).
		Return(errorx.Internal("failed to persist register results"))

	cd := &service{
		lambdaRepo: mockLambRepo,
		oauthRepo:  mockOAuthRepo,
		amazon:     mockAmazon,
	}

	register, err := cd.Register(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, errorx.Internal("failed to persist register results"), err)
	assert.Nil(t, register)
}

func TestRemoveSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLambRepo := testdata.NewMockLambda(ctrl)
	mockAmazon := testdata.NewMockAmazon(ctrl)
	mockOAuthRepo := testdata.NewMockOAuth(ctrl)

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimSub.Str(), "account_name")
	request := &dto.ReqURILambda{
		Lambda: "file1",
	}
	functionARN := "arn:aws:lambda:us-east-2:123456789012:function:file1"
	scheduleARN := "arn:aws:scheduler:us-east-2:123456789012:schedule/default/file1"
	accountID := uint64(123)

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, "account_name").Times(1).
		Return(&dto.RespUser{ID: accountID}, nil)

	mockLambRepo.EXPECT().LambdaInfo(ctx, accountID, "file1").Times(1).
		Return(&dto.RespInfo{
			FunctionName: "org_name-account_name-file1",
			FunctionArn:  functionARN,
			Scheduler: dto.Scheduler{
				ScheduleArn:  scheduleARN,
				ScheduleName: "org_name-account_name-file1",
			},
		}, nil)

	mockAmazon.EXPECT().RemoveLambda(ctx, gomock.Any()).Times(1).
		DoAndReturn(func(_ *gin.Context, input *lambda.DeleteFunctionInput) (*lambda.DeleteFunctionOutput, error) {
			assert.Equal(t, "org_name-account_name-file1", *input.FunctionName)
			return &lambda.DeleteFunctionOutput{
				ResultMetadata: func() middleware.Metadata {
					m := middleware.Metadata{}
					m.Set("val", 200)
					return m
				}(),
			}, nil
		})

	mockAmazon.EXPECT().RemoveScheduler(ctx, gomock.Any()).Times(1).
		DoAndReturn(func(_ *gin.Context, input *scheduler.DeleteScheduleInput) (*scheduler.DeleteScheduleOutput, error) {
			assert.Equal(t, "org_name-account_name-file1", *input.Name)
			return &scheduler.DeleteScheduleOutput{
				ResultMetadata: func() middleware.Metadata {
					m := middleware.Metadata{}
					m.Set("val", 200)
					return m
				}(),
			}, nil
		})

	mockLambRepo.EXPECT().DeleteLambdaTX(ctx, gomock.Any(), gomock.Any()).Times(1).
		Return(nil)

	cd := &service{
		lambdaRepo: mockLambRepo,
		amazon:     mockAmazon,
		oauthRepo:  mockOAuthRepo,
	}

	remove, err := cd.Remove(ctx, request)
	assert.NoError(t, err)
	assert.Equal(t, &dto.RespRemove{
		Lambdas: dto.RespLamBrief{
			Name: "org_name-account_name-file1",
			Arn:  functionARN,
		},
		Scheduler: dto.RespSchBrief{
			Arn: scheduleARN,
		},
	}, remove)
}

func TestRemoveUserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOAuthRepo := testdata.NewMockOAuth(ctrl)

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimSub.Str(), "account_name")
	request := &dto.ReqURILambda{
		Lambda: "file1",
	}

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, "account_name").Times(1).
		Return(nil, errorx.NotFound("user not found"))

	cd := &service{
		oauthRepo: mockOAuthRepo,
	}

	remove, err := cd.Remove(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, errorx.NotFound("user not found"), err)
	assert.Nil(t, remove)
}

func TestRemoveLambdaNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLambRepo := testdata.NewMockLambda(ctrl)
	mockOAuthRepo := testdata.NewMockOAuth(ctrl)
	ctx := new(gin.Context)
	ctx.Set(constant.ClaimSub.Str(), "account_name")
	request := &dto.ReqURILambda{
		Lambda: "file1",
	}
	accountID := uint64(123)

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, "account_name").Times(1).
		Return(&dto.RespUser{ID: accountID}, nil)

	mockLambRepo.EXPECT().LambdaInfo(ctx, accountID, "file1").Times(1).
		Return(nil, errorx.NotFound("lambda not found"))

	cd := &service{
		lambdaRepo: mockLambRepo,
		oauthRepo:  mockOAuthRepo,
	}

	remove, err := cd.Remove(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, errorx.NotFound("lambda not found"), err)
	assert.Nil(t, remove)
}

func TestRemoveLambdaError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLambRepo := testdata.NewMockLambda(ctrl)
	mockAmazon := testdata.NewMockAmazon(ctrl)
	mockOAuthRepo := testdata.NewMockOAuth(ctrl)
	ctx := new(gin.Context)
	ctx.Set(constant.ClaimSub.Str(), "account_name")
	request := &dto.ReqURILambda{
		Lambda: "file1",
	}
	functionARN := "arn:aws:lambda:us-east-2:123456789012:function:file1"
	accountID := uint64(123)

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, "account_name").Times(1).
		Return(&dto.RespUser{ID: accountID}, nil)

	mockLambRepo.EXPECT().LambdaInfo(ctx, accountID, "file1").Times(1).
		Return(&dto.RespInfo{
			FunctionArn:  functionARN,
			FunctionName: "org_name-account_name-file1",
		}, nil)

	mockAmazon.EXPECT().RemoveLambda(ctx, gomock.Any()).Times(1).
		Return(nil, errorx.Internal("failed to remove lambda"))

	cd := &service{
		lambdaRepo: mockLambRepo,
		amazon:     mockAmazon,
		oauthRepo:  mockOAuthRepo,
	}

	remove, err := cd.Remove(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, errorx.Internal("failed to remove lambda"), err)
	assert.Nil(t, remove)
}

func TestRemoveSchedulerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLambRepo := testdata.NewMockLambda(ctrl)
	mockAmazon := testdata.NewMockAmazon(ctrl)
	mockOAuthRepo := testdata.NewMockOAuth(ctrl)

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimSub.Str(), "account_name")
	request := &dto.ReqURILambda{
		Lambda: "file1",
	}
	functionARN := "arn:aws:lambda:us-east-2:123456789012:function:file1"
	scheduleARN := "arn:aws:scheduler:us-east-2:123456789012:schedule/default/file1"
	accountID := uint64(123)

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, "account_name").Times(1).
		Return(&dto.RespUser{ID: accountID}, nil)

	mockLambRepo.EXPECT().LambdaInfo(ctx, accountID, "file1").Times(1).
		Return(&dto.RespInfo{
			FunctionArn:  functionARN,
			FunctionName: "org_name-account_name-file1",
			Scheduler: dto.Scheduler{
				ScheduleArn:  scheduleARN,
				ScheduleName: "org_name-account_name-file1",
			},
		}, nil)

	mockAmazon.EXPECT().RemoveLambda(ctx, gomock.Any()).Times(1).
		Return(&lambda.DeleteFunctionOutput{
			ResultMetadata: func() middleware.Metadata {
				m := middleware.Metadata{}
				m.Set("val", 200)
				return m
			}(),
		}, nil)

	mockAmazon.EXPECT().RemoveScheduler(ctx, gomock.Any()).Times(1).
		Return(nil, errorx.Internal("failed to remove scheduler"))

	cd := &service{
		lambdaRepo: mockLambRepo,
		amazon:     mockAmazon,
		oauthRepo:  mockOAuthRepo,
	}

	remove, err := cd.Remove(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, errorx.Internal("failed to remove scheduler"), err)
	assert.Nil(t, remove)
}

func TestRemoveDeleteLambdaTXError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLambRepo := testdata.NewMockLambda(ctrl)
	mockAmazon := testdata.NewMockAmazon(ctrl)
	mockOAuthRepo := testdata.NewMockOAuth(ctrl)

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimSub.Str(), "account_name")
	request := &dto.ReqURILambda{
		Lambda: "file1",
	}
	functionARN := "arn:aws:lambda:us-east-2:123456789012:function:file1"
	accountID := uint64(123)

	mockOAuthRepo.EXPECT().FindUserByAcn(ctx, "account_name").Times(1).
		Return(&dto.RespUser{ID: accountID}, nil)

	mockLambRepo.EXPECT().LambdaInfo(ctx, accountID, "file1").Times(1).
		Return(&dto.RespInfo{
			FunctionArn:  functionARN,
			FunctionName: "org_name-account_name-file1",
		}, nil)

	mockAmazon.EXPECT().RemoveLambda(ctx, gomock.Any()).Times(1).
		Return(&lambda.DeleteFunctionOutput{
			ResultMetadata: func() middleware.Metadata {
				m := middleware.Metadata{}
				m.Set("val", 200)
				return m
			}(),
		}, nil)

	mockLambRepo.EXPECT().DeleteLambdaTX(ctx, gomock.Any(), gomock.Any()).Times(1).
		Return(errorx.Internal("failed to delete lambda tx"))

	cd := &service{
		lambdaRepo: mockLambRepo,
		amazon:     mockAmazon,
		oauthRepo:  mockOAuthRepo,
	}

	remove, err := cd.Remove(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, errorx.Internal("failed to delete lambda tx"), err)
	assert.Nil(t, remove)
}

func TestLogSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedMessage := `{"timestamp":1714857600000,"message":"log-message"}`
	mockAmazon := testdata.NewMockAmazon(ctrl)

	mockAmazon.EXPECT().DescribeLogStreams(gomock.Any(), gomock.Any()).AnyTimes().Return(&cloudwatchlogs.DescribeLogStreamsOutput{
		LogStreams: []types.LogStream{
			{
				LogStreamName: aws.String("log-stream-name"),
			},
		},
	}, nil)

	mockAmazon.EXPECT().GetLogEvents(gomock.Any(), gomock.Any()).AnyTimes().Return(&cloudwatchlogs.GetLogEventsOutput{
		Events: []types.OutputLogEvent{
			{
				Timestamp: aws.Int64(1714857600000),
				Message:   aws.String(expectedMessage),
			},
		},
	}, nil)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create a new gin.Context
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = r

		// create WebSocket upgrader
		upgrader := &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		request := &dto.ReqURILambda{
			Lambda: "file1",
		}

		// call Logs method
		svc := &service{
			amazon: mockAmazon,
		}
		err := svc.Logs(ctx, request, upgrader)
		if err != nil {
			t.Fatalf("Logs method failed: %v", err)
		}
	}))
	defer server.Close()

	// convert HTTP URL to WebSocket URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// create WebSocket client connection
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer conn.Close()

	// read the message sent from the server
	_, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	assert.Equal(t, expectedMessage, string(message))
}

func TestLogsWSUpgradeError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req, err := http.NewRequest("GET", "/logs", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	svc := &service{}

	err = svc.Logs(ctx, nil, upgrader)
	assert.Error(t, err)
	assert.Equal(t, `failed to upgrade websocket: websocket: the client is not using the websocket protocol: 'upgrade' token not found in 'Connection' header`, err.Error())
}

func TestLogsDescribeLogStreamsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAmazon := testdata.NewMockAmazon(ctrl)

	mockAmazon.EXPECT().DescribeLogStreams(gomock.Any(), gomock.Any()).AnyTimes().
		Return(nil, errorx.Internal("failed to describe log streams"))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = r

		upgrader := &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		request := &dto.ReqURILambda{
			Lambda: "file1",
		}

		svc := &service{
			amazon: mockAmazon,
		}
		err := svc.Logs(ctx, request, upgrader)
		assert.Error(t, err)
		assert.Equal(t, errorx.Internal("failed to describe log streams"), err)
	}))
	defer server.Close()
}

func TestLogsGetLogEventsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAmazon := testdata.NewMockAmazon(ctrl)

	mockAmazon.EXPECT().DescribeLogStreams(gomock.Any(), gomock.Any()).AnyTimes().Return(&cloudwatchlogs.DescribeLogStreamsOutput{
		LogStreams: []types.LogStream{
			{
				LogStreamName: aws.String("log-stream-name"),
			},
		},
	}, nil)

	mockAmazon.EXPECT().GetLogEvents(gomock.Any(), gomock.Any()).AnyTimes().
		Return(nil, errorx.Internal("failed to get log events"))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = r

		upgrader := &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		request := &dto.ReqURILambda{
			Lambda: "file1",
		}

		svc := &service{
			amazon: mockAmazon,
		}
		err := svc.Logs(ctx, request, upgrader)
		assert.Error(t, err)
		assert.Equal(t, errorx.Internal("failed to get log events"), err)
	}))
	defer server.Close()
}

func TestLogsLogStreamNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAmazon := testdata.NewMockAmazon(ctrl)

	mockAmazon.EXPECT().DescribeLogStreams(gomock.Any(), gomock.Any()).AnyTimes().Return(&cloudwatchlogs.DescribeLogStreamsOutput{
		LogStreams: []types.LogStream{},
	}, nil)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = r

		upgrader := &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		request := &dto.ReqURILambda{
			Lambda: "file1",
		}

		svc := &service{
			amazon: mockAmazon,
		}
		err := svc.Logs(ctx, request, upgrader)
		assert.Error(t, err)
		assert.Equal(t, errorx.NotFound("no log streams found"), err)
	}))
	defer server.Close()
}
