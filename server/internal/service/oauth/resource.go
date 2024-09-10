package oauth

import (
	"net/http"

	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/57blocks/auto-action/server/internal/third-party/logx"

	"github.com/gin-gonic/gin"
)

var ServiceImpl Service

func Login(c *gin.Context) {
	req := new(dto.ReqLogin)

	if err := c.BindJSON(req); err != nil {
		c.Error(errorx.BadRequest(err.Error()))
		return
	}

	resp, err := ServiceImpl.Login(c, *req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func Logout(c *gin.Context) {
	req := new(dto.ReqLogout)

	if err := c.BindJSON(req); err != nil {
		c.Error(errorx.BadRequest(err.Error()))
		return
	}

	logx.Logger.DEBUG("oauth logout", map[string]interface{}{
		"logout_token": req.Token,
	})

	resp, err := ServiceImpl.Logout(c, *req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func Refresh(c *gin.Context) {
	req := new(dto.ReqRefresh)

	if err := c.BindJSON(req); err != nil {
		c.Error(errorx.BadRequest(err.Error()))
		return
	}

	logx.Logger.DEBUG("oauth refresh", map[string]interface{}{
		"refresh_token": req.Refresh,
	})

	resp, err := ServiceImpl.Refresh(c, *req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
