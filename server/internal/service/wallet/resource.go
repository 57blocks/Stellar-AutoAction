package wallet

import (
	"net/http"

	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

var ServiceImpl Service

func Create(c *gin.Context) {
	resp, err := ServiceImpl.Create(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, resp)
}

func Remove(c *gin.Context) {
	req := new(dto.RemoveWalletReqInfo)

	if err := c.BindUri(req); err != nil {
		c.Error(errorx.BadRequest(err.Error()))
		c.Abort()
		return
	}

	err := ServiceImpl.Remove(c, req)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}

func ListWallets(c *gin.Context) {
	resp, err := ServiceImpl.List(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, resp)
}

func Verify(c *gin.Context) {
	req := new(dto.VerifyWalletReqInfo)

	if err := c.BindUri(req); err != nil {
		c.Error(errorx.BadRequest(err.Error()))
		return
	}

	resp, err := ServiceImpl.Verify(c, req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
