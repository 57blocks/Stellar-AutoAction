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

//go:generate mockgen -destination ./oauth_mock.go -package repo -source oauth.go OAuth
type (
	OAuth interface {
		FindUserByAcn(c context.Context, acn string) (*dto.RespUser, error)
		FindUserByOrgAcn(c context.Context, req *dto.ReqOrgAcn) (*dto.RespUser, error)

		FindOrg(c context.Context, id uint64) (*dto.RespOrg, error)
		FindOrgByName(c context.Context, name string) (*dto.RespOrg, error)

		FindTokenByRefreshID(c context.Context, refresh string) (*model.Token, error)
		SyncToken(c context.Context, token *model.Token) error
		DeleteTokenByAccess(c context.Context, access string) error
	}

	oauth struct {
		Instance *db.Instance
	}
)

var OAuthRepo OAuth

func NewOAuth() {
	if OAuthRepo == nil {
		OAuthRepo = &oauth{
			Instance: db.Inst,
		}
	}
}

func (o *oauth) FindUserByAcn(c context.Context, acn string) (*dto.RespUser, error) {
	u := new(dto.RespUser)
	if err := db.Inst.Conn(c).Table(model.TabNamUser()).
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

func (o *oauth) FindUserByOrgAcn(c context.Context, req *dto.ReqOrgAcn) (*dto.RespUser, error) {
	u := new(dto.RespUser)
	if err := o.Instance.Conn(c).Table(model.TabNamUserAbbr()).
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

func (o *oauth) FindOrg(c context.Context, id uint64) (*dto.RespOrg, error) {
	org := new(dto.RespOrg)
	if err := o.Instance.Conn(c).Table(model.TabNamOrg()).
		Where(map[string]interface{}{
			"id": id,
		}).
		First(org).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NotFound("organization not found")
		}

		return nil, errorx.Internal(err.Error())
	}

	return org, nil
}

func (o *oauth) FindOrgByName(c context.Context, name string) (*dto.RespOrg, error) {
	org := new(dto.RespOrg)
	if err := o.Instance.Conn(c).Table(model.TabNamOrg()).
		Where(map[string]interface{}{
			"name": name,
		}).
		First(org).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NotFound("organization not found")
		}

		return nil, errorx.Internal(err.Error())
	}

	return org, nil
}

func (o *oauth) FindTokenByRefreshID(c context.Context, refreshID string) (*model.Token, error) {
	t := new(model.Token)
	if err := o.Instance.Conn(c).Table(t.TableName()).
		Where(map[string]interface{}{
			"refresh_id": refreshID,
		}).
		First(t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NotFound("none refresh token found")
		}

		return nil, errorx.Internal(err.Error())
	}

	return t, nil
}

func (o *oauth) SyncToken(c context.Context, token *model.Token) error {
	if err := o.Instance.Conn(c).
		Table(token.TableName()).
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

func (o *oauth) DeleteTokenByAccess(c context.Context, access string) error {
	if err := o.Instance.Conn(c).
		Table(model.TabNamToken()).
		Where(map[string]interface{}{
			"access": access,
		}).
		Delete(&model.Token{}).Error; err != nil {
		return errorx.Internal(err.Error())
	}

	return nil
}
