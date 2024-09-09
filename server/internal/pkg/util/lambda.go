package util

import (
	"context"
	"fmt"

	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	svcOrg "github.com/57blocks/auto-action/server/internal/service/organization"

	"github.com/gin-gonic/gin"
)

func GenLambdaFuncName(c context.Context, name string) string {
	org, _ := svcOrg.Conductor.Organization(c)

	ctx := c.(*gin.Context)
	account, _ := ctx.Get("jwt_account")

	return fmt.Sprintf("%s-%s-%s", org.Name, account, name)
}

func GenEventPayload(c context.Context) (*dto.StdEventPayload, error) {
	ctx, ok := c.(*gin.Context)
	if !ok {
		return nil, errorx.GinContextConv()
	}

	jwtOrg, _ := ctx.Get("jwt_organization")
	jwtAccount, _ := ctx.Get("jwt_account")

	return &dto.StdEventPayload{
		Organization: jwtOrg.(string),
		Account:      jwtAccount.(string),
	}, nil
}
