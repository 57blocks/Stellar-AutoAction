package organization

import (
	"net/http"

	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	dto "github.com/57blocks/auto-action/server/internal/service/dto/organization"

	"github.com/gin-gonic/gin"
)

func Keys(c *gin.Context) {
	req := new(dto.ReqKeys)

	if err := c.ShouldBindJSON(req); err != nil {
		c.Error(errorx.BadRequest(err.Error()))
		return
	}

	roleKey, err := Conductor.OrgRoleKey(c, req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, roleKey.CSRoleKeys)
}
