package lambda

import (
	"net/http"

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
