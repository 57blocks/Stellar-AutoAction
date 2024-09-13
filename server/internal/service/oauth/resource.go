package oauth

import (
	"net/http"

	"github.com/57blocks/auto-action/server/internal/constant"
	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

var ServiceImpl Service

func Login(c *gin.Context) {
	req := new(dto.ReqLogin)

	if err := c.BindJSON(req); err != nil {
		c.Error(errorx.BadRequest(err.Error()))
		c.Abort()
		return
	}

	resp, err := ServiceImpl.Login(c, *req)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, resp)
}

func Logout(c *gin.Context) {
	raw, ok := c.Get(constant.ClaimRaw.Str())
	if !ok {
		c.Error(errorx.Unauthorized())
		c.Abort()
		return
	}
	rawStr, ok := raw.(string)
	if !ok {
		c.Error(errorx.Internal("raw is not string"))
		c.Abort()
		return
	}

	resp, err := ServiceImpl.Logout(c, rawStr)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, resp)
}

func Refresh(c *gin.Context) {
	raw, ok := c.Get(constant.ClaimRaw.Str())
	if !ok {
		c.Error(errorx.Unauthorized())
		c.Abort()
		return
	}
	rawStr, ok := raw.(string)
	if !ok {
		c.Error(errorx.Internal("raw is not string"))
		c.Abort()
		return
	}

	resp, err := ServiceImpl.Refresh(c, rawStr)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, resp)
}
