package oauth

import (
	"net/http"

	"github.com/57blocks/auto-action/server/internal/constant"
	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

type (
	Resource interface {
		Signup(c *gin.Context)
		Login(c *gin.Context)
		Logout(c *gin.Context)
		Refresh(c *gin.Context)
	}
	resource struct {
		service OAuthService
	}
)

var ResourceImpl Resource

func NewOAuthResource() {
	if ResourceImpl == nil {
		ResourceImpl = &resource{
			service: OAuthServiceImpl,
		}
	}
}

func (re *resource) Signup(c *gin.Context) {
	req := new(dto.ReqSignup)

	if err := c.BindJSON(req); err != nil {
		c.Error(errorx.BadRequest(err.Error()))
		c.Abort()
		return
	}

	if err := re.service.Signup(c, *req); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (re *resource) Login(c *gin.Context) {
	req := new(dto.ReqLogin)

	if err := c.BindJSON(req); err != nil {
		c.Error(errorx.BadRequest(err.Error()))
		c.Abort()
		return
	}

	resp, err := re.service.Login(c, *req)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (re *resource) Logout(c *gin.Context) {
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

	resp, err := re.service.Logout(c, rawStr)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (re *resource) Refresh(c *gin.Context) {
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

	resp, err := re.service.Refresh(c, rawStr)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, resp)
}
