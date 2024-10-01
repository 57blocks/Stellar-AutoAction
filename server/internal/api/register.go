package api

import (
	"net/http"

	"github.com/57blocks/auto-action/server/internal/api/middleware"
	"github.com/57blocks/auto-action/server/internal/service/lambda"
	"github.com/57blocks/auto-action/server/internal/service/oauth"
	"github.com/57blocks/auto-action/server/internal/service/wallet"

	"github.com/gin-gonic/gin"
)

func RegisterHandlers(g *gin.Engine) http.Handler {
	oauthGroup := g.Group("/oauth")
	{
		oauthGroup.POST("/signup", oauth.ResourceImpl.Signup)
		oauthGroup.POST("/login", oauth.ResourceImpl.Login)
		oauthGroup.DELETE("/logout", middleware.AuthHeader(), oauth.ResourceImpl.Logout)
		oauthGroup.POST("/refresh", middleware.AuthHeader(), oauth.ResourceImpl.Refresh)
	}

	lambdaGroup := g.Group("/lambda", middleware.Authentication(), middleware.Authorization())
	{
		lambdaGroup.POST("", middleware.RegisterESLintCheck(), lambda.ResourceImpl.Register)
		lambdaGroup.POST("/:lambda", lambda.ResourceImpl.Invoke)
		lambdaGroup.GET("", lambda.ResourceImpl.List)
		lambdaGroup.GET("/:lambda", lambda.ResourceImpl.Info)
		lambdaGroup.GET("/:lambda/logs", lambda.ResourceImpl.Logs)
		lambdaGroup.DELETE("/:lambda", lambda.ResourceImpl.Remove)
	}

	walletGroup := g.Group("/wallet", middleware.Authentication(), middleware.Authorization())
	{
		walletGroup.GET("", wallet.List)
		walletGroup.POST("", wallet.Create)
		walletGroup.DELETE("/:address", wallet.Remove)
		walletGroup.POST("/:address", wallet.Verify)
	}

	return g
}
