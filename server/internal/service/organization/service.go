package organization

import (
	"context"
	"errors"
	"fmt"

	"github.com/57blocks/auto-action/server/internal/db"
	"github.com/57blocks/auto-action/server/internal/model/oauth"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	Service interface {
		Organization(c context.Context) (*oauth.Organization, error)
	}
	conductor struct{}
)

var Conductor Service

func init() {
	if Conductor == nil {
		Conductor = &conductor{}
	}
}

func (cd conductor) Organization(c context.Context) (*oauth.Organization, error) {
	ctx, ok := c.(*gin.Context)
	if !ok {
		return nil, errorx.GinContextConv()
	}

	jwtOrg, _ := ctx.Get("jwt_organization")

	org := new(oauth.Organization)
	if err := db.Conn(c).Table(org.TableName()).
		Where(map[string]interface{}{
			"name": jwtOrg,
		}).
		First(org).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NotFound("none organization found")
		}

		return nil, errorx.Internal(fmt.Sprintf("find organization by name: %s, occurred error: %s", jwtOrg, err.Error()))
	}

	return org, nil
}
