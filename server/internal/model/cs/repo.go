package cs

import (
	"context"

	"github.com/57blocks/auto-action/server/internal/db"
	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/model/oauth"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"

	"gorm.io/gorm"
)

//go:generate mockgen -destination ./repo_mock.go -package cs -source repo.go Repo
type (
	Repo interface {
		ToSign(c context.Context, req *dto.ReqToSign) ([]*dto.RespToSign, error)
	}
	conductor struct{}
)

var Conductor Repo

func init() {
	if Conductor == nil {
	}
	Conductor = &conductor{}
}
func (cs *conductor) ToSign(c context.Context, req *dto.ReqToSign) ([]*dto.RespToSign, error) {
	roles := make([]*dto.RespToSign, 0)
	if err := db.Conn(c).
		Table(TabNamCSRoleAbbr()).
		Joins("LEFT JOIN organization AS o ON o.id = csr.organization_id").
		Joins("LEFT JOIN \"user\" AS u ON u.id = csr.account_id").
		Preload("Organization", func(db *gorm.DB) *gorm.DB {
			return db.Table(oauth.TabNamOrg())
		}).
		Preload("Account", func(db *gorm.DB) *gorm.DB {
			return db.Table(oauth.TabNamUser())
		}).
		Preload("Keys", func(db *gorm.DB) *gorm.DB {
			return db.Table(TabNameCSKey())
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
