package oauth

import (
	"net/http"

	"github.com/57blocks/auto-action/server/internal/api/middleware"
	pkgLog "github.com/57blocks/auto-action/server/internal/pkg/log"

	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {
	pkgLog.Logger.DEBUG(
		middleware.NewFromCTX(ctx).Format(),
		map[string]any{
			"account": "v3nooom@outlook.com",
		},
		map[string]any{
			"account2": "v3nooom@outlook.com2",
			"account3": "v3nooom@outlook.com3",
		},
	)

	pkgLog.Logger.DEBUG(
		middleware.NewFromCTX(ctx).Format(),
	)

	pkgLog.Logger.INFO(middleware.NewFromCTX(ctx).Format())

	ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
}
