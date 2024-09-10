package repo

import (
	"context"
	"errors"

	"github.com/57blocks/auto-action/server/internal/config"
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
		ToSign(c context.Context, userID uint64, from string) (*dto.RespCSKey, error)

		SyncCSKey(c context.Context, key *model.CubeSignerKey) error
	}
	cubeSigner struct {
		Instance *db.Instance
	}
)

var CubeSignerRepo CubeSigner

func NewCubeSigner() {
	if CubeSignerRepo == nil {
		CubeSignerRepo = &cubeSigner{
			Instance: db.Inst,
		}
	}
}
func (cs *cubeSigner) ToSign(c context.Context, userID uint64, from string) (*dto.RespCSKey, error) {
	csKey := new(dto.RespCSKey)
	if err := cs.Instance.Conn(c).
		Table(model.TabNameCSKey()).
		Where(map[string]interface{}{
			"account_id": userID,
			"key":        from,
		}).
		Preload("Account", func(db *gorm.DB) *gorm.DB {
			return db.Table(model.TabNamUser())
		}).
		First(csKey).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NotFound("cube signer key not found")
		}
		return nil, errorx.Internal(err.Error())
	}

	csKey.Organization = config.GlobalConfig.CS.Organization

	return csKey, nil
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
