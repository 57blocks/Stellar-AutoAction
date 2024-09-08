package oauth

import (
	"context"
	"errors"

	"github.com/57blocks/auto-action/server/internal/db"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate mockgen -destination ./repo_mock.go -package oauth -source repo.go Repo
type (
	Repo interface {
		FindUserByAcn(c context.Context, acn string) (*RespUser, error)
		FindUserByOrgAcn(c context.Context, req ReqOrgAcn) (*RespUser, error)

		FindOrg(c context.Context, id uint64) (*RespOrg, error)
		FindOrgByName(c context.Context, name string) (*RespOrg, error)

		SyncToken(c context.Context, token *Token) error
	}

	conductor struct{}
)

var Conductor Repo

func init() {
	if Conductor == nil {
	}
	Conductor = &conductor{}
}

func (cd *conductor) FindUserByAcn(c context.Context, acn string) (*RespUser, error) {
	u := new(RespUser)
	if err := db.Conn(c).Table(TabNamUser()).
		Where(map[string]interface{}{
			"account": acn,
		}).
		First(u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NotFound("user/organization not found")
		}

		return nil, errorx.Internal(err.Error())
	}

	return u, nil
}

func (cd *conductor) FindUserByOrgAcn(c context.Context, req ReqOrgAcn) (*RespUser, error) {
	u := new(RespUser)
	if err := db.Conn(c).Table(TabNamUserAbbr()).
		Joins("LEFT JOIN organization AS o ON u.organization_id = o.id").
		Where(map[string]interface{}{
			"u.account": req.AcnName,
			"o.name":    req.OrgName,
		}).
		First(u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NotFound("user/organization not found")
		}

		return nil, errorx.Internal(err.Error())
	}

	return u, nil
}

func (cd *conductor) FindOrg(c context.Context, id uint64) (*RespOrg, error) {
	o := new(RespOrg)
	if err := db.Conn(c).Table(TabNamOrg()).
		Where(map[string]interface{}{
			"id": id,
		}).
		First(o).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NotFound("organization not found")
		}

		return nil, errorx.Internal(err.Error())
	}

	return o, nil
}

func (cd *conductor) FindOrgByName(c context.Context, name string) (*RespOrg, error) {
	o := new(RespOrg)
	if err := db.Conn(c).Table(TabNamOrg()).
		Where(map[string]interface{}{
			"name": name,
		}).
		First(o).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NotFound("organization not found")
		}

		return nil, errorx.Internal(err.Error())
	}

	return o, nil
}

func (cd *conductor) SyncToken(c context.Context, token *Token) error {
	if err := db.Conn(c).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			UpdateAll: true,
		}).
		Create(token).
		Error; err != nil {
		return errorx.Internal(err.Error())
	}

	return nil
}
