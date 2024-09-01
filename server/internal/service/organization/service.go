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
		CurrentOrg(c context.Context) (*model.Organization, error)
		CurrentOrgKPs(c context.Context) (*dto.RespRelatedKeyPairs, error)
		CurrentVpc(c context.Context) (*model.Vpc, error)
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

func (cd conductor) CurrentOrg(c context.Context) (*model.Organization, error) {
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

func (cd conductor) CurrentOrgKPs(c context.Context) (*dto.RespRelatedKeyPairs, error) {
	ctx, ok := c.(*gin.Context)
	if !ok {
		return nil, errors.New("convert context.Context to gin.Context failed")
	}

	jwtOrg, _ := ctx.Get("jwt_organization")

	okpList := make([]*model.OrganizationKeyPairs, 0)
	if err := db.Conn(c).
		Table(new(model.OrganizationKeyPairs).TableNameWithAbbr()).
		Joins("LEFT JOIN organization AS o ON o.id = okp.organization_id").
		Where(map[string]interface{}{"o.name": jwtOrg}).
		Find(&okpList).Error; err != nil {
		return nil, errors.Wrap(err, "query organization key pairs failed")
	}

	if len(okpList) == 0 {
		return nil, errors.New("none organization related key pairs found")
	}

	cubeKeyPairs := make([]dto.CubeSignerPairs, 0, len(okpList))
	for _, okp := range okpList {
		cubeKeyPairs = append(cubeKeyPairs, dto.CubeSignerPairs{
			Private: okp.PrivateKey,
			Public:  okp.PublicKey,
		})
	}

	return &dto.RespRelatedKeyPairs{
		JWTPairs:        dto.JWTPairs{},
		CubeSignerPairs: cubeKeyPairs,
	}, nil
}

func (cd conductor) CurrentVpc(c context.Context) (*model.Vpc, error) {
	ctx, ok := c.(*gin.Context)
	if !ok {
		return nil, errors.New("convert context.Context to gin.Context failed")
	}

	jwtOrg, _ := ctx.Get("jwt_organization")

	vpc := new(model.Vpc)
	if err := db.Conn(c).Table(vpc.TableNameWithAbbr()).
		Joins("LEFT JOIN organization AS o ON o.id = v.organization_id").
		Where(map[string]interface{}{
			"o.name": jwtOrg,
		}).
		First(vpc).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("none vpc found")
		}

		return nil, errors.Wrap(err, "db error when find vpc by organization name")
	}

	return vpc, nil
}
