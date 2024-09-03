package lambda

import (
	"context"
	"fmt"

	dtoLam "github.com/57blocks/auto-action/server/internal/service/dto/lambda"
	svcOrg "github.com/57blocks/auto-action/server/internal/service/organization"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func genLambdaFuncName(c context.Context, name string) string {
	org, _ := svcOrg.Conductor.Organization(c)

	ctx := c.(*gin.Context)
	account, _ := ctx.Get("jwt_account")

	return fmt.Sprintf("%s-%s-%s", org.Name, account, name)
}

func genSchEventPayload(c context.Context) (dtoLam.ReqSchedulerEvent, error) {
	// TODO: make key-value pairs extraction as an util
	ctx, ok := c.(*gin.Context)
	if !ok {
		return dtoLam.ReqSchedulerEvent{}, errors.New("convert context.Context to gin.Context failed")
	}

	jwtOrg, _ := ctx.Get("jwt_organization")
	jwtAccount, _ := ctx.Get("jwt_account")

	return dtoLam.ReqSchedulerEvent{
		APIKey:       "sample_api_key", // TODO
		Organization: jwtOrg.(string),
		Account:      jwtAccount.(string),
	}, nil
}
