package repo

import (
	"context"

	"github.com/57blocks/auto-action/server/internal/db"
	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/model"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"

	"gorm.io/gorm"
)

//go:generate mockgen -destination ./cs_mock.go -package repo -source cs.go CubeSigner
type (
	CubeSigner interface {
		ToSign(c context.Context, req *dto.ReqToSign) ([]*dto.RespToSign, error)
	}
	cubeSigner struct {
		Instance *db.Instance
	}
)

var CubeSignerImpl CubeSigner

func NewCubeSigner() {
	if CubeSignerImpl == nil {
		CubeSignerImpl = &cubeSigner{
			Instance: db.Inst,
		}
	}
}
func (cs *cubeSigner) ToSign(c context.Context, req *dto.ReqToSign) ([]*dto.RespToSign, error) {
	roles := make([]*dto.RespToSign, 0)
	if err := cs.Instance.Conn(c).
		Table(model.TabNamCSRoleAbbr()).
		Joins("LEFT JOIN organization AS o ON o.id = csr.organization_id").
		Joins("LEFT JOIN \"user\" AS u ON u.id = csr.account_id").
		Preload("Organization", func(db *gorm.DB) *gorm.DB {
			return db.Table(model.TabNamOrg())
		}).
		Preload("Account", func(db *gorm.DB) *gorm.DB {
			return db.Table(model.TabNamUser())
		}).
		Preload("Keys", func(db *gorm.DB) *gorm.DB {
			return db.Table(model.TabNameCSKey())
		}).
		Where(map[string]interface{}{
			"o.name":    req.Organization,
			"u.account": req.Account,
		}).
		Find(&roles).Error; err != nil {
		return nil, errorx.Internal(err.Error())
	}

	if len(roles) == 0 {
		return nil, errorx.NotFound("none roles found to sign")
	}

	return roles, nil
}
