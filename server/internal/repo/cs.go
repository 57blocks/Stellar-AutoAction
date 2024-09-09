package repo

import (
	"context"
	"errors"

	"github.com/57blocks/auto-action/server/internal/db"
	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/model"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate mockgen -destination ./cs_mock.go -package repo -source cs.go CubeSigner
type (
	CubeSigner interface {
		ToSign(c context.Context, req *dto.ReqToSign) ([]*dto.RespToSign, error)

		FindCSByOrgAcn(c context.Context, req *dto.ReqCSRole) (*dto.RespCSRole, error)
		SyncCSKey(c context.Context, key *model.CubeSignerKey) error
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

func (cs *cubeSigner) FindCSByOrgAcn(c context.Context, req *dto.ReqCSRole) (*dto.RespCSRole, error) {
	role := new(dto.RespCSRole)
	if err := cs.Instance.Conn(c).
		Table(model.TabNamCSRole()).
		Where(map[string]interface{}{
			"organization_id": req.OrgID,
			"account_id":      req.AcnID,
		}).
		First(role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NotFound("role not found")
		}
		return nil, errorx.Internal(err.Error())
	}

	return role, nil
}

func (cs *cubeSigner) SyncCSKey(c context.Context, key *model.CubeSignerKey) error {
	if err := cs.Instance.Conn(c).
		Table(key.TableName()).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "key"}},
			UpdateAll: true,
		}).
		Create(key).Error; err != nil {
		return errorx.Internal(err.Error())
	}

	return nil
}
