package middleware

import (
	"fmt"

	pkgLog "github.com/57blocks/auto-action/server/internal/pkg/log"

	"github.com/gin-gonic/gin"
)

func NewFromCTX(ctx *gin.Context) *RestFormatter {
	return &RestFormatter{
		Method: ctx.Request.Method,
		Path:   ctx.Request.URL.Path,
		Status: ctx.Writer.Status(),
	}
}

type RestFormatter struct {
	Method string
	Path   string
	Status int
}

func (rf *RestFormatter) Format() string {
	return fmt.Sprintf(
		"%s %s - %v",
		rf.Method,
		rf.Path,
		rf.Status,
	)
}

type (
	Formatter interface {
		Format() string
	}
	Logger interface {
		DEBUG(msg string, argMap map[string]interface{})
		INFO(msg string, argMap map[string]interface{})
		WARN(msg string, argMap map[string]interface{})
		ERROR(msg string, argMap map[string]interface{})
	}
)

func ZapLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		pkgLog.Logger.DEBUG(NewFromCTX(c).Format())
	}
}
