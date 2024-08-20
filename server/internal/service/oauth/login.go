package oauth

import (
	"encoding/json"
	"io"
	"net/http"

	model "github.com/57blocks/auto-action/server/internal/dto/oauth"
	"github.com/57blocks/auto-action/server/internal/pkg/jwtx"
	pkgLog "github.com/57blocks/auto-action/server/internal/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Login(ctx *gin.Context) {
	req := new(model.Request)

	jsonData, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := json.Unmarshal(jsonData, req); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	pkgLog.Logger.DEBUG("login", map[string]interface{}{
		"account":      req.Account,
		"organization": req.Organization,
		"environment":  req.Environment,
		"password":     req.Password,
	})

	tokens, err := jwtx.Assign()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	resp := model.BuildResp(
		model.WithAccount(req.Account),
		model.WithOrganization(req.Organization),
		model.WithBound(model.BuildBound(
			model.WithBoundName(req.Environment),
			model.WithBoundEndPoint(viper.GetString("bound.endpoint")),
		)),
		model.WithTokens(tokens),
	)

	ctx.JSON(http.StatusOK, resp)
}
