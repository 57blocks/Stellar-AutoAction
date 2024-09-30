package util

import (
	"context"
	"testing"

	"github.com/57blocks/auto-action/server/internal/constant"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDecodeBase64StringSuccess(t *testing.T) {
	encodedStr := "dGVzdA=="

	ctx := context.Background()

	decodedStr, err := DecodeBase64String(ctx, &encodedStr)

	assert.NoError(t, err)
	assert.Equal(t, "test", decodedStr)
}

func TestDecodeBase64StringFailed(t *testing.T) {
	encodedStr := "invalid"

	ctx := context.Background()

	decodedStr, err := DecodeBase64String(ctx, &encodedStr)

	assert.Error(t, err)
	assert.Equal(t, "", decodedStr)
}

func TestGetRoleNameSuccess(t *testing.T) {
	org := "test-org"
	account := "test-account"

	ctx := context.Background()

	roleName := GetRoleName(ctx, org, account)

	assert.Equal(t, "AA-test-org-test-account-Role", roleName)
}

func TestGetSecretNameSuccess(t *testing.T) {
	org := "test-org"
	account := "test-account"

	ctx := context.Background()

	secretName := GetSecretName(ctx, org, account)

	assert.Equal(t, "AA_test-org_test-account_SEC", secretName)
}

func TestGenLambdaFuncNameSuccess(t *testing.T) {
	org := "test-org"
	account := "test-account"
	name := "test-function"

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimIss.Str(), org)
	ctx.Set(constant.ClaimSub.Str(), account)

	lambdaFuncName := GenLambdaFuncName(ctx, name)

	assert.Equal(t, "test-org-test-account-test-function", lambdaFuncName)
}

func TestGetInputPayloadSuccess(t *testing.T) {
	payload := `{"key": "value"}`

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimIss.Str(), "test-org")
	ctx.Set(constant.ClaimSub.Str(), "test-account")

	inputPayload, err := GenEventPayload(ctx, payload)

	assert.NoError(t, err)
	assert.Equal(t, "test-org", (*inputPayload)["organization"])
	assert.Equal(t, "test-account", (*inputPayload)["account"])
	assert.Equal(t, "value", (*inputPayload)["key"])
}

func TestGetInputPayloadFailed(t *testing.T) {
	payload := `invalid json`

	ctx := new(gin.Context)
	ctx.Set(constant.ClaimIss.Str(), "test-org")
	ctx.Set(constant.ClaimSub.Str(), "test-account")

	inputPayload, err := GenEventPayload(ctx, payload)

	assert.Error(t, err)
	assert.Equal(t, errorx.Internal("failed to unmarshal payload: invalid character 'i' looking for beginning of value"), err)
	assert.Nil(t, inputPayload)
}
