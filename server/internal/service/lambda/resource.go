package lambda

import (
	"net/http"

	dto "github.com/57blocks/auto-action/server/internal/service/dto/lambda"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	r := c.Request

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp, err := Conductor.Register(c, r)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func Trigger(c *gin.Context) {
	req := new(dto.ReqTrigger)

	if err := c.BindUri(req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp, err := Conductor.Trigger(c, req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func Info(c *gin.Context) {
	req := new(dto.ReqInfo)

	if err := c.BindUri(req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp, err := Conductor.Info(c, req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func Logs(c *gin.Context) {
	req := new(dto.ReqLogs)

	if err := c.BindUri(req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := Conductor.Logs(c, req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
