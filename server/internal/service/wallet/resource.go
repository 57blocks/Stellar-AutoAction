package wallet

import (
	"net/http"

	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/gin-gonic/gin"
)

var ServiceImpl Service

func Create(c *gin.Context) {
	r := c.Request

	resp, err := ServiceImpl.Create(c, r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func Remove(c *gin.Context) {
	req := new(dto.RemoveWalletReqInfo)

	if err := c.BindUri(req); err != nil {
		c.Error(errorx.BadRequest(err.Error()))
		return
	}

	err := ServiceImpl.Remove(c, req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func List(c *gin.Context) {
	r := c.Request

	resp, err := ServiceImpl.List(c, r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
