package wallet

import (
	"net/http"

	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

type (
	Resource interface {
		Create(c *gin.Context)
		Remove(c *gin.Context)
		List(c *gin.Context)
		Verify(c *gin.Context)
	}
	resource struct {
		service WalletService
	}
)

var ResourceImpl Resource

func NewWalletResource() {
	if ResourceImpl == nil {
		ResourceImpl = &resource{
			service: WalletServiceImpl,
		}
	}
}

func (re *resource) Create(c *gin.Context) {
	resp, err := re.service.Create(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (re *resource) Remove(c *gin.Context) {
	req := new(dto.ReqRemoveWallet)

	if err := c.BindUri(req); err != nil {
		c.Error(errorx.BadRequest(err.Error()))
		c.Abort()
		return
	}

	err := re.service.Remove(c, req)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (re *resource) List(c *gin.Context) {
	resp, err := re.service.List(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (re *resource) Verify(c *gin.Context) {
	req := new(dto.ReqVerifyWallet)

	if err := c.BindUri(req); err != nil {
		c.Error(errorx.BadRequest(err.Error()))
		return
	}

	resp, err := re.service.Verify(c, req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
