package wallet

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Create(c *gin.Context) {
	r := c.Request

	resp, err := Conductor.Create(c, r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
