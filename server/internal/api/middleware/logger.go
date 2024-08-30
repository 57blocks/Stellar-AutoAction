package middleware

import (
	"fmt"

	pkgLog "github.com/57blocks/auto-action/server/internal/pkg/log"

	"github.com/gin-gonic/gin"
)

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

type FormatterREST struct {
	Method string
	Path   string
	Status int
}

func (rf *FormatterREST) Format() string {
	return fmt.Sprintf(
		"%s %s - %v",
		rf.Method,
		rf.Path,
		rf.Status,
	)
}

func NewFromContext(ctx *gin.Context) *FormatterREST {
	return &FormatterREST{
		Method: ctx.Request.Method,
		Path:   ctx.Request.URL.Path,
		Status: ctx.Writer.Status(),
	}
}

func ZapLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		pkgLog.Logger.DEBUG(NewFromContext(c).Format())
		c.Next()
	}
}
