package api

import (
	"net/http"

	"github.com/57blocks/auto-action/server/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

func Boot() *gin.Engine {
	g := gin.New(
		WithCustomRecovery(),
		WithReqHeader(),
		WithError(),
		WithResponse(),
	)

	g.Use(middleware.ZapLogger())

	g.GET("/up", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	RegisterHandlers(g)

	return g
}
