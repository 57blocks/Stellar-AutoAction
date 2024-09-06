package lambda

import (
	"context"
	"fmt"

	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	dtoLam "github.com/57blocks/auto-action/server/internal/service/dto/lambda"
	svcOrg "github.com/57blocks/auto-action/server/internal/service/organization"

	"github.com/gin-gonic/gin"
)

func genLambdaFuncName(c context.Context, name string) string {
	org, _ := svcOrg.Conductor.Organization(c)

	ctx := c.(*gin.Context)
	account, _ := ctx.Get("jwt_account")

	return fmt.Sprintf("%s-%s-%s", org.Name, account, name)
}

func genEventPayload(c context.Context) (*dtoLam.StdEventPayload, error) {
	ctx, ok := c.(*gin.Context)
	if !ok {
		return nil, errorx.GinContextConv()
	}

	jwtOrg, _ := ctx.Get("jwt_organization")
	jwtAccount, _ := ctx.Get("jwt_account")

	return &dtoLam.StdEventPayload{
		Organization: jwtOrg.(string),
		Account:      jwtAccount.(string),
	}, nil
}
