package api

import (
	"net/http"

	"github.com/57blocks/auto-action/server/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

func WithCustomRecovery() gin.OptionFunc {
	return func(g *gin.Engine) {
		g.Use(
			gin.CustomRecovery(func(c *gin.Context, err interface{}) {
				e, ok := err.(error)
				if ok {
					c.JSON(http.StatusInternalServerError, gin.H{
						"success": false,
						"message": "internal error",
						"msg":     e.Error(),
					})
				}
			},
			),
		)
	}
}

func WithCORS() gin.OptionFunc {
	return func(g *gin.Engine) {
		g.Use(middleware.CORS())
	}
}

func WithReqHeader() gin.OptionFunc {
	return func(g *gin.Engine) {
		g.Use(middleware.Header())
	}
}

func WithResponse() gin.OptionFunc {
	return func(g *gin.Engine) {
		g.Use(middleware.Response())
	}
}

func WithError() gin.OptionFunc {
	return func(g *gin.Engine) {
		g.Use(middleware.Error())
	}
}
