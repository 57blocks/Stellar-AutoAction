package cs

import (
	"errors"
	"os"
	"testing"

	"github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/third-party/amazonx"
	"github.com/57blocks/auto-action/server/internal/third-party/logx"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// Before test, setup log
func TestMain(m *testing.M) {
	testConfig := config.Configuration{
		Log: config.Log{
			Level:    "debug",
			Encoding: "json",
		},
	}
	logx.Setup(&testConfig)

	os.Exit(m.Run())
}

func TestCubeSignerTokenSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockAmazon := amazonx.NewMockAmazon(ctrl)
	mockAmazon.EXPECT().GetSecretValue(gomock.Any(), gomock.Any()).
		Return(&secretsmanager.GetSecretValueOutput{
			SecretString: aws.String(`{"token": "test-token"}`),
		}, nil)

	svc := &service{
		amazon: mockAmazon,
	}

	token, err := svc.CubeSignerToken(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "test-token", token)
}

func TestCubeSignerTokenInternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockAmazon := amazonx.NewMockAmazon(ctrl)
	mockAmazon.EXPECT().GetSecretValue(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("get secret value error"))

	svc := &service{
		amazon: mockAmazon,
	}
	_, err := svc.CubeSignerToken(ctx)
	assert.Error(t, err)
	assert.Equal(t, "get secret value error", err.Error())
}

func TestCubeSignerTokenJsonError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockAmazon := amazonx.NewMockAmazon(ctrl)
	mockAmazon.EXPECT().GetSecretValue(gomock.Any(), gomock.Any()).
		Return(&secretsmanager.GetSecretValueOutput{
			SecretString: aws.String(`"unmarshal error"`),
		}, nil)

	svc := &service{
		amazon: mockAmazon,
	}
	_, err := svc.CubeSignerToken(ctx)
	assert.Error(t, err)
	assert.Equal(t, "json unmarshal error when parse secret value: json: cannot unmarshal string into Go value of type map[string]interface {}", err.Error())
}

func TestGetSecRoleSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockAmazon := amazonx.NewMockAmazon(ctrl)
	mockAmazon.EXPECT().GetSecretValue(gomock.Any(), gomock.Any()).
		Return(&secretsmanager.GetSecretValueOutput{
			SecretString: aws.String(`{"cs_role": "Role#1234"}`),
		}, nil)

	svc := &service{
		amazon: mockAmazon,
	}

	role, err := svc.GetSecRole(ctx, "test-secret")
	assert.NoError(t, err)
	assert.Equal(t, "Role#1234", role)
}

func TestGetSecRoleInternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockAmazon := amazonx.NewMockAmazon(ctrl)
	mockAmazon.EXPECT().GetSecretValue(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("get secret value error"))

	svc := &service{
		amazon: mockAmazon,
	}

	_, err := svc.GetSecRole(ctx, "test-secret")
	assert.Error(t, err)
	assert.Equal(t, "get secret value error", err.Error())
}

func TestGetSecRoleJsonError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := new(gin.Context)
	mockAmazon := amazonx.NewMockAmazon(ctrl)
	mockAmazon.EXPECT().GetSecretValue(gomock.Any(), gomock.Any()).
		Return(&secretsmanager.GetSecretValueOutput{
			SecretString: aws.String(`"unmarshal error"`),
		}, nil)

	svc := &service{
		amazon: mockAmazon,
	}
	_, err := svc.GetSecRole(ctx, "test-secret")
	assert.Error(t, err)
	assert.Equal(t, "json unmarshal error when parse secret value: json: cannot unmarshal string into Go value of type map[string]interface {}", err.Error())
}
