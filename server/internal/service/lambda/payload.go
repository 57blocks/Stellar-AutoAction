package lambda

import (
	"context"

	"github.com/57blocks/auto-action/server/internal/config"
	dtoLam "github.com/57blocks/auto-action/server/internal/service/dto/lambda"
	dtoOrg "github.com/57blocks/auto-action/server/internal/service/dto/organization"
	svcOrg "github.com/57blocks/auto-action/server/internal/service/organization"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func buildLambdaEvent(c context.Context) (dtoLam.ReqSchedulerEvent, error) {
	// TODO: make key-value pairs extraction as an util
	ctx, ok := c.(*gin.Context)
	if !ok {
		return dtoLam.ReqSchedulerEvent{}, errors.New("convert context.Context to gin.Context failed")
	}

	jwtRaw, _ := ctx.Get("jwt_raw")

	keyPairs, err := svcOrg.Conductor.CurrentOrgKPs(c)
	if err != nil {
		return dtoLam.ReqSchedulerEvent{}, err
	}

	return dtoLam.ReqSchedulerEvent{
		Token: jwtRaw.(string),
		JWTPairs: dtoOrg.JWTPairs{
			Private: config.Global.JWT.PrivateKey,
			Public:  config.Global.JWT.PublicKey,
		},
		CubeSignerPairs: keyPairs.CubeSignerPairs,
	}, nil
}
