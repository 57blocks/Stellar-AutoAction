package lambda

import (
	"net/http"

	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	dto "github.com/57blocks/auto-action/server/internal/service/dto/lambda"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	r := c.Request

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		c.Error(errorx.BadRequest(err.Error()))
		return
	}

	resp, err := Conductor.Register(c, r)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func Invoke(c *gin.Context) {
	req := new(dto.ReqInvoke)

	if err := c.BindUri(req); err != nil {
		c.Error(errorx.BadRequest(err.Error()))
		return
	}

	if err := c.ShouldBindJSON(req); err != nil {
		c.Error(errorx.BadRequest(err.Error()))
		return
	}

	resp, err := Conductor.Invoke(c, req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func Info(c *gin.Context) {
	req := new(dto.ReqInfo)

	if err := c.BindUri(req); err != nil {
		c.Error(errorx.BadRequest(err.Error()))
		return
	}

	resp, err := Conductor.Info(c, req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func Logs(c *gin.Context) {
	req := new(dto.ReqLogs)

	if err := c.BindUri(req); err != nil {
		c.Error(errorx.BadRequest(err.Error()))
		return
	}

	if err := Conductor.Logs(c, req); err != nil {
		c.Error(err)
		return
	}
}
