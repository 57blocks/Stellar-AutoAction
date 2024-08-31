package oauth

import (
	"net/http"

	pkgLog "github.com/57blocks/auto-action/server/internal/pkg/log"
	oauth2 "github.com/57blocks/auto-action/server/internal/service/dto/oauth"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	req := new(oauth2.ReqLogin)

	if err := c.BindJSON(req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp, err := Conductor.Login(c, *req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func Logout(c *gin.Context) {
	req := new(oauth2.ReqLogout)

	if err := c.BindJSON(req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pkgLog.Logger.DEBUG("logout", map[string]interface{}{
		"logout_token": req.Token,
	})

	resp, err := Conductor.Logout(c, *req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func Refresh(c *gin.Context) {
	req := new(oauth2.ReqRefresh)

	if err := c.BindJSON(req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pkgLog.Logger.DEBUG("refresh", map[string]interface{}{
		"refresh_token": req.Refresh,
	})

	resp, err := Conductor.Refresh(c, *req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
