package cs

import (
	"net/http"

	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

var ServiceImpl Service

func ToSign(c *gin.Context) {
	req := new(dto.ReqToSign)

	if err := c.ShouldBindJSON(req); err != nil {
		c.Error(errorx.BadRequest(err.Error()))
		return
	}

	respToSign, err := ServiceImpl.ToSign(c, req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, respToSign)
}
