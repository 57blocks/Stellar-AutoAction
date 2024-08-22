package oauth

import (
	"context"
	"time"

	"github.com/57blocks/auto-action/server/internal/model"
	"github.com/57blocks/auto-action/server/internal/pkg/db"
	"github.com/57blocks/auto-action/server/internal/pkg/jwtx"
	dto "github.com/57blocks/auto-action/server/internal/service/dto/oauth"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	Service interface {
		Login(c context.Context, req dto.ReqLogin) (*dto.RespCredential, error)
		Refresh(c context.Context, req dto.ReqRefresh) (*dto.RespCredential, error)
		Logout(c context.Context, req dto.ReqLogout) (*dto.RespLogout, error)
	}
	ServiceConductor struct{}
)

var Conductor Service

func init() {
	if Conductor == nil {
		Conductor = &ServiceConductor{}
	}
}

func (sc *ServiceConductor) Login(c context.Context, req dto.ReqLogin) (*dto.RespCredential, error) {
	user := new(model.User)
	if err := db.Conn(c).Table(user.TableNameWithAbbr()).
		Joins("LEFT JOIN principal_organization AS po ON pu.organization_id = po.id").
		Where(map[string]interface{}{
			"pu.account": req.Account,
			"po.name":    req.Organization,
		}).
		First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user/organization not found")
		}
		return nil, errors.New(err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), req.Password); err != nil {
		return nil, errors.New("password not match")
	}

	// tokens assignment
	now := time.Now().UTC()
	tokenExp := now.AddDate(0, 1, 0)
	refreshExp := now.AddDate(0, 3, 0)

	tc := jwtx.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "v3nooom", // TODO
			IssuedAt:  time.Now().UTC().Unix(),
			Subject:   "st3llar",
			ExpiresAt: tokenExp.Unix(),
		},
		Account:      "Account_sample",
		Organization: "Organization_sample",
	}

	rf := jwt.StandardClaims{
		Issuer:    "v3nooom",
		IssuedAt:  time.Now().UTC().Unix(),
		Subject:   "st3llar",
		ExpiresAt: refreshExp.Unix(),
	}

	tokens, err := jwtx.Assign(jwtx.ClaimPair{
		Token:   tc,
		Refresh: rf,
	})
	if err != nil {
		return nil, err
	}

	// sync token pairs
	token := &model.Token{
		Access:         tokens.Token,
		Refresh:        tokens.Refresh,
		UserId:         user.ID,
		AccessExpires:  tokenExp,
		RefreshExpires: refreshExp,
	}
	if err := db.Conn(c).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			UpdateAll: true,
		}).
		Create(token).
		Error; err != nil {
		return nil, errors.New(err.Error())
	}

	// build response
	resp := dto.BuildRespCred(
		dto.WithAccount(req.Account),
		dto.WithOrganization(req.Organization),
		dto.WithEnvironment(req.Environment),
		dto.WithTokenPair(tokens),
	)

	return resp, nil
}

func (sc *ServiceConductor) Refresh(c context.Context, req dto.ReqRefresh) (*dto.RespCredential, error) {
	token := new(model.Token)
	if err := db.Conn(c).
		Where(map[string]interface{}{
			"refresh":             req.Refresh,
			"refresh_expires > ?": time.Now().UTC(),
		}).
		First(token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("refresh token not found")
		}
		return nil, errors.New(err.Error())
	}

	refreshExp := token.RefreshExpires.UTC()
	if refreshExp.Before(time.Now().UTC()) {
		return nil, errors.New("refresh token expired")
	}

	now := time.Now().UTC()
	tokenExp := now.AddDate(0, 1, 0)

	tc := jwtx.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "v3nooom", // TODO
			IssuedAt:  time.Now().UTC().Unix(),
			Subject:   "st3llar",
			ExpiresAt: tokenExp.Unix(),
		},
		Account:      "Account_sample",
		Organization: "Organization_sample",
	}

	rf := jwt.StandardClaims{
		Issuer:    "v3nooom",
		IssuedAt:  time.Now().UTC().Unix(),
		Subject:   "st3llar",
		ExpiresAt: refreshExp.Unix(),
	}

	tokens, err := jwtx.Assign(jwtx.ClaimPair{
		Token:   tc,
		Refresh: rf,
	})
	if err != nil {
		return nil, err
	}

	// save tokens association
	token := &model.Token{
		Access:         tokens.Token,
		Refresh:        tokens.Refresh,
		UserId:         user.ID,
		AccessExpires:  tokenExp,
		RefreshExpires: refreshExp,
	}
	if err := db.Conn(c).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			UpdateAll: true,
		}).
		Create(token).
		Error; err != nil {
		return nil, errors.New(err.Error())
	}

	resp := dto.BuildRespCred(
		dto.WithAccount(req.Account),
		dto.WithOrganization(req.Organization),
		dto.WithEnvironment(req.Environment),
		dto.WithTokenPair(tokens),
	)

	return resp, nil
}

func (sc *ServiceConductor) Logout(c context.Context, req dto.ReqLogout) (*dto.RespLogout, error) {
	if err := db.Conn(c).
		Where(map[string]interface{}{
			"access": req.Token,
		}).
		Delete(&model.Token{}).Error; err != nil {
		return nil, errors.New(err.Error())
	}

	return new(dto.RespLogout), nil
}
