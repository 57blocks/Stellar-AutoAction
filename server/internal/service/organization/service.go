package organization

import (
	"context"

	"github.com/57blocks/auto-action/server/internal/pkg/db"
	dto "github.com/57blocks/auto-action/server/internal/service/dto/organization"
	"github.com/57blocks/auto-action/server/internal/service/model"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type (
	Service interface {
		RelatedOrg(c context.Context) (*model.Organization, error)
		RelatedKeyPairs(c context.Context) (*dto.RespRelatedKeyPairs, error)
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

func (cd conductor) RelatedOrg(c context.Context) (*model.Organization, error) {
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
		return nil, errors.Wrap(err, "none organization failed")
	}

	return org, nil
}

func (cd conductor) RelatedKeyPairs(c context.Context) (*dto.RespRelatedKeyPairs, error) {
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
