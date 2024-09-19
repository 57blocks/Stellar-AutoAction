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
		oauthGroup.POST("/login", oauth.Login)
		oauthGroup.DELETE("/logout", middleware.AuthHeader(), oauth.Logout)
		oauthGroup.POST("/refresh", middleware.AuthHeader(), oauth.Refresh)
	}

	lambdaGroup := g.Group("/lambda", middleware.Authentication(), middleware.Authorization())
	{
		lambdaGroup.POST("", middleware.RegisterESLintCheck(), lambda.Register)
		lambdaGroup.POST("/:lambda", lambda.Invoke)
		lambdaGroup.GET("", lambda.List)
		lambdaGroup.GET("/:lambda", lambda.Info)
		lambdaGroup.GET("/:lambda/logs", lambda.Logs)
		lambdaGroup.DELETE("/:lambda", lambda.Remove)
	}

	walletGroup := g.Group("/wallet", middleware.Authentication(), middleware.Authorization())
	{
		walletGroup.GET("", wallet.ListWallets)
		walletGroup.POST("", wallet.Create)
		walletGroup.DELETE("/:address", wallet.Remove)
		walletGroup.POST("/:address", wallet.Verify)
	}

	return g
}
