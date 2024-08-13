package oauth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Refresh(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
}
