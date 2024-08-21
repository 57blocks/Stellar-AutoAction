package oauth

import (
	"encoding/json"
	"io"
	"net/http"

	pkgLog "github.com/57blocks/auto-action/server/internal/pkg/log"
	dto "github.com/57blocks/auto-action/server/internal/service/dto/oauth"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	req := dto.Request{}

	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := json.Unmarshal(jsonData, &req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	pkgLog.Logger.DEBUG("login", map[string]interface{}{
		"account":      req.Account,
		"organization": req.Organization,
		"environment":  req.Environment,
		"password":     string(req.Password),
	})

	resp, err := Conductor.Login(c, req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func Refresh(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
