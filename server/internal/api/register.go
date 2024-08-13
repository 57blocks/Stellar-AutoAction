package api

import (
	"net/http"

	"github.com/57blocks/auto-action/server/internal/api/middleware"
	svcLambda "github.com/57blocks/auto-action/server/internal/service/lambda"
	svcAuth "github.com/57blocks/auto-action/server/internal/service/oauth"

	"github.com/gin-gonic/gin"
)

func RegisterHandlers(g *gin.Engine) http.Handler {
	oauth := g.Group("/oauth")
	{
		oauth.POST("/login", svcAuth.Login)

		oauth.POST("/logout", svcAuth.Logout)
		oauth.POST("/refresh", svcAuth.Refresh)
	}

	lambda := g.Group("/lambda", middleware.Authentication(), middleware.Authorization())
	{
		lambda.POST("/register", svcLambda.Register)
	}

	return g
}
