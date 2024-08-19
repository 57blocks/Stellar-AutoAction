package oauth

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"net/http"

	"github.com/57blocks/auto-action/server/internal/pkg/jwtx"
	pkgLog "github.com/57blocks/auto-action/server/internal/pkg/log"

	"github.com/gin-gonic/gin"
)

type (
	ReqLogin struct {
		_            struct{}
		Account      string `json:"account"`
		Organization string `json:"organization"`
		Password     []byte `json:"password"`
		Environment  string `json:"environment"`
	}

	RespLogin struct {
		_            struct{}
		Account      string `json:"account" toml:"account"`
		Organization string `json:"organization" toml:"organization"`
		*jwtx.Tokens `json:"tokens" toml:"tokens"`
		*Bound       `json:"bound" toml:"bound"`
	}
	RespLoginOpt func(cred *RespLogin)

	Bound struct {
		_        struct{}
		Name     string `json:"name" toml:"name"`
		EndPoint string `json:"endpoint" toml:"endpoint"`
	}
	RespBoundOpt func(bound *Bound)
)

func Login(ctx *gin.Context) {
	req := new(ReqLogin)

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

	resp := BuildResp(
		WithAccount(req.Account),
		WithOrganization(req.Organization),
		WithBound(BuildBound(
			WithBoundName(req.Environment),
			WithBoundEndPoint(viper.GetString("bound.endpoint")),
		)),
		WithTokens(tokens),
	)
	fmt.Printf("%#v\n", resp)

	ctx.JSON(http.StatusOK, resp)
}
