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
		WithLogger(),
		WithPostError(),
		WithPostResponse(),
	)

	GinEngine.GET("/up", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	RegisterHandlers(GinEngine)

	return nil
}

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

func WithLogger() gin.OptionFunc {
	return func(g *gin.Engine) {
		g.Use(middleware.ZapLogger())
	}
}

func WithPostResponse() gin.OptionFunc {
	return func(g *gin.Engine) {
		g.Use(middleware.PostResponse())
	}
}

func WithPostError() gin.OptionFunc {
	return func(g *gin.Engine) {
		g.Use(middleware.PostHandleErr())
	}
}
