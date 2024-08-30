package api

import (
	"net/http"

	"github.com/57blocks/auto-action/server/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

var GinEngine *gin.Engine

func Setup() error {
	GinEngine = gin.New(
		WithCustomRecovery(),
		WithReqHeader(),
		WithError(),
		WithResponse(),
	)

	GinEngine.Use(
		middleware.ZapLogger(),
		middleware.Header(),
	)

	GinEngine.GET("/up", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	RegisterHandlers(GinEngine)

	return nil
}
