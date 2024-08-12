package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Boot() *gin.Engine {
	//gin.LoggerWithConfig(gin.LoggerConfig{})

	g := gin.New(
		WithCustomRecovery(),
		WithReqHeader(),
		WithError(),
		WithResponse(),
	)

	g.GET("/up", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	RegisterHandlers(g)

	return g
}
