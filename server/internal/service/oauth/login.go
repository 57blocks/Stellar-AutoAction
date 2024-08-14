package oauth

import (
	"encoding/json"
	"io"
	"net/http"

	pkgLog "github.com/57blocks/auto-action/server/internal/pkg/log"
	"github.com/57blocks/auto-action/server/internal/third_party/jwtx"

	"github.com/gin-gonic/gin"
)

type (
	Req struct {
		_            struct{}
		Account      string `json:"account"`
		Organization string `json:"organization"`
		Password     []byte `json:"password"`
		Environment  string `json:"environment"`
	}
	Resp struct {
		_            struct{}
		Account      string `json:"account" toml:"account"`
		Organization string `json:"organization" toml:"organization"`
		*jwtx.Tokens `json:"tokens" toml:"tokens"`
		*Environment `json:"environment" toml:"environment"`
	}
	CredOpt func(cred *Resp)

	Environment struct {
		_        struct{}
		Name     string `json:"name" toml:"name"`
		EndPoint string `json:"endpoint" toml:"endpoint"`
	}
)

func Login(ctx *gin.Context) {
	req := new(Req)
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

	resp := &Resp{
		Account:      req.Account,
		Organization: req.Organization,
		Environment: &Environment{
			Name:     req.Environment,
			EndPoint: "http://st3llar-alb-365211.us-east-2.elb.amazonaws.com/",
		},
		Tokens: tokens,
	}

	ctx.JSON(http.StatusOK, resp)
}
