package util

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/57blocks/auto-action/server/internal/constant"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

func GenLambdaFuncName(c context.Context, name string) string {
	ctx := c.(*gin.Context)
	org, _ := ctx.Get(constant.ClaimIss.Str())
	account, _ := ctx.Get(constant.ClaimSub.Str())

	return fmt.Sprintf("%s-%s-%s", org, account, name)
}

func GenEventPayload(c context.Context, payload string) (*map[string]interface{}, error) {
	ctx, ok := c.(*gin.Context)
	if !ok {
		return nil, errorx.GinContextConv()
	}

	jwtOrg, _ := ctx.Get(constant.ClaimIss.Str())
	jwtAccount, _ := ctx.Get(constant.ClaimSub.Str())

	inputPayload := make(map[string]interface{})
	if len(payload) > 0 {
		json.Unmarshal([]byte(payload), &inputPayload) // the payload is validated in CLI already
	}

	inputPayload["organization"] = jwtOrg.(string)
	inputPayload["account"] = jwtAccount.(string)

	return &inputPayload, nil
}

func GetRoleName(c context.Context, org, account string) string {
	return fmt.Sprintf("AA-%s-%s-Role", org, account)
}

func GetSecretName(c context.Context, org string, account string) string {
	return fmt.Sprintf("AA_%s_%s_SEC", org, account)
}

func DecodeBase64String(encodedStr *string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(*encodedStr)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}
