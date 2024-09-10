package wallet

import (
	"net/http"

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
