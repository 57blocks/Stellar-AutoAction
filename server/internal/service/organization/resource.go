package organization

import (
	"net/http"

	dto "github.com/57blocks/auto-action/server/internal/service/dto/organization"

	"github.com/gin-gonic/gin"
)

func Keys(c *gin.Context) {
	req := new(dto.ReqKeys)

	if err := c.ShouldBindJSON(req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	roleKey, err := Conductor.OrgRoleKey(c, req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, roleKey.CSRoleKeys)
}
