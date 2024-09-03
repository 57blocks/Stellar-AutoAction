package organization

import (
	"context"

	"github.com/57blocks/auto-action/server/internal/db"
	"github.com/57blocks/auto-action/server/internal/model"
	dto "github.com/57blocks/auto-action/server/internal/service/dto/organization"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type (
	Service interface {
		Organization(c context.Context) (*model.Organization, error)
		OrgRoleKey(c context.Context) (*dto.RespRelatedRoleKey, error)
		OrgSecret(c context.Context) (*model.OrgSecret, error)
	}
	conductor struct{}
)

var (
	Conductor Service
)

func init() {
	if Conductor == nil {
		Conductor = &conductor{}
	}
}

func (cd conductor) Organization(c context.Context) (*model.Organization, error) {
	ctx, ok := c.(*gin.Context)
	if !ok {
		return nil, errors.New("convert context.Context to gin.Context failed")
	}

	jwtOrg, _ := ctx.Get("jwt_organization")

	org := new(model.Organization)
	if err := db.Conn(c).Table(org.TableName()).
		Where(map[string]interface{}{
			"name": jwtOrg,
		}).
		First(org).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("none organization found")
		}

		return nil, errors.Wrap(err, "db error when find organization by name")
	}

	return org, nil
}

func (cd conductor) OrgRoleKey(c context.Context) (*dto.RespRelatedRoleKey, error) {
	ctx, ok := c.(*gin.Context)
	if !ok {
		return nil, errors.New("convert context.Context to gin.Context failed")
	}

	jwtOrg, _ := ctx.Get("jwt_organization")

	orkList := make([]*model.CSOrgRoleKey, 0)
	if err := db.Conn(c).
		Table(model.TabNameAbbrOrgRoleKey()).
		Joins("LEFT JOIN organization AS o ON o.id = ork.organization_id").
		Where(map[string]interface{}{
			"o.name": jwtOrg,
		}).
		Find(&orkList).Error; err != nil {
		return nil, errors.Wrap(err, "query organization key pairs failed")
	}

	if len(orkList) == 0 {
		return nil, errors.New("none organization related key pairs found")
	}

	csRoleKeyList := make([]dto.RespCSRoleKey, 0, len(orkList))
	for _, ork := range orkList {
		csRoleKeyList = append(csRoleKeyList, dto.RespCSRoleKey{
			CSRoleID: ork.CSRoleID,
			CSKeyID:  ork.CSKeyID,
			CSScopes: ork.CSScopes,
		})
	}

	return &dto.RespRelatedRoleKey{
		CSRoleKeys: csRoleKeyList,
	}, nil
}

func (cd conductor) OrgSecret(c context.Context) (*model.OrgSecret, error) {
	ctx, ok := c.(*gin.Context)
	if !ok {
		return nil, errors.New("convert context.Context to gin.Context failed")
	}

	jwtOrg, _ := ctx.Get("jwt_organization")

	orgSecret := new(model.OrgSecret)
	if err := db.Conn(c).Table(model.TabNameAbbrOrgSecret()).
		Joins("LEFT JOIN organization AS o ON o.id = os.organization_id").
		Where(map[string]interface{}{
			"o.name":    jwtOrg,
			"os.active": true,
		}).
		First(orgSecret).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("none organization secret found")
		}

		return nil, errors.Wrap(err, "db error when find organization secret by organization id")
	}

	return orgSecret, nil
}
