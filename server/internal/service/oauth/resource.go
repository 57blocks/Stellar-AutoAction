package oauth

import (
	"net/http"

	pkgLog "github.com/57blocks/auto-action/server/internal/pkg/log"
	dto "github.com/57blocks/auto-action/server/internal/service/dto/oauth"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	req := new(dto.ReqLogin)

	if err := c.BindJSON(req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pkgLog.Logger.DEBUG("login", map[string]interface{}{
		"account":      req.Account,
		"organization": req.Organization,
		"environment":  req.Environment,
		"password":     string(req.Password),
	})

	resp, err := Conductor.Login(c, *req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func Logout(c *gin.Context) {
	req := new(dto.ReqLogout)

	if err := c.BindJSON(req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pkgLog.Logger.DEBUG("logout", map[string]interface{}{
		"token": req.Token,
	})

	resp, err := Conductor.Logout(c, *req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func Refresh(c *gin.Context) {
	req := new(dto.ReqRefresh)

	if err := c.BindJSON(req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pkgLog.Logger.DEBUG("logout", map[string]interface{}{
		"refresh": req.Refresh,
	})

	resp, err := Conductor.Refresh(c, *req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
