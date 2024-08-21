package api

import (
	"net/http"

	"github.com/57blocks/auto-action/server/internal/api/middleware"
	"github.com/57blocks/auto-action/server/internal/service/lambda"
	"github.com/57blocks/auto-action/server/internal/service/oauth"

	"github.com/gin-gonic/gin"
)

func RegisterHandlers(g *gin.Engine) http.Handler {
	oauthGroup := g.Group("/oauth")
	{
		oauthGroup.POST("/login", oauth.Login)
		oauthGroup.POST("/logout", oauth.Logout)
		oauthGroup.POST("/refresh", oauth.Refresh)
	}

	lambdaGroup := g.Group("/lambda", middleware.Authentication(), middleware.Authorization())
	{
		lambdaGroup.POST("/register", lambda.Register)
	}

	return g
}
