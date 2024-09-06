package organization

import (
	"fmt"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"net/http"

	dto "github.com/57blocks/auto-action/server/internal/service/dto/organization"

	"github.com/gin-gonic/gin"
)

func Keys(c *gin.Context) {
	req := new(dto.ReqKeys)

	if err := c.ShouldBindJSON(req); err != nil {
		c.Error(fmt.Errorf("%w", errorx.BadRequestErr(err.Error())))
		return
	}

	roleKey, err := Conductor.OrgRoleKey(c, req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, roleKey.CSRoleKeys)
}
