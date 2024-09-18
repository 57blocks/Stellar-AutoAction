package util

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/57blocks/auto-action/server/internal/constant"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	svcOrg "github.com/57blocks/auto-action/server/internal/service/organization"
	"github.com/57blocks/auto-action/server/internal/third-party/amazonx"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/gin-gonic/gin"
)

func GenLambdaFuncName(c context.Context, name string) string {
	org, _ := svcOrg.ServiceImpl.Organization(c)

	ctx := c.(*gin.Context)
	account, _ := ctx.Get(constant.ClaimSub.Str())

	return fmt.Sprintf("%s-%s-%s", org.Name, account, name)
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

func GetRoleARN(c context.Context, roleName string) (string, error) {
	input := &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	}

	result, err := amazonx.Conductor.GetRole(c, input)
	if err != nil {
		return "", errorx.Internal(fmt.Sprintf("get role arn error: %v", err))
	}

	return *result.Role.Arn, nil
}
