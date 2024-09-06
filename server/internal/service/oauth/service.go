package oauth

import (
	"context"
	"errors"
	"time"

	"github.com/57blocks/auto-action/server/internal/db"
	"github.com/57blocks/auto-action/server/internal/model"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/57blocks/auto-action/server/internal/pkg/jwtx"
	dto "github.com/57blocks/auto-action/server/internal/service/dto/oauth"

	"github.com/dgrijalva/jwt-go"
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
	conductor struct{}
)

var Conductor Service

func init() {
	if Conductor == nil {
		Conductor = &conductor{}
	}
}

func (cd *conductor) Login(c context.Context, req dto.ReqLogin) (*dto.RespCredential, error) {
	user := new(model.User)
	if err := db.Conn(c).Table(user.TableNameWithAbbr()).
		Joins("LEFT JOIN organization AS o ON u.organization_id = o.id").
		Where(map[string]interface{}{
			"u.account": req.Account,
			"o.name":    req.Organization,
		}).
		First(user).Error; err != nil {
		if errors.As(err, &gorm.ErrRecordNotFound) {
			return nil, errorx.NotFound("user/organization not found")
		}
		return nil, errorx.Internal(err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), req.Password); err != nil {
		return nil, errorx.BadRequest("password not match")
	}

	// tokens assignment
	now := time.Now().UTC()
	tokenExp := now.AddDate(0, 1, 0)
	refreshExp := now.AddDate(0, 3, 0)

	access, err := jwtx.AssignAccess(jwtx.AccessClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "v3nooom", // TODO
			IssuedAt:  time.Now().UTC().Unix(),
			Subject:   "st3llar",
			ExpiresAt: tokenExp.Unix(),
		},
		Account:      req.Account,
		Organization: req.Organization,
		Environment:  req.Environment,
	})
	if err != nil {
		return nil, err
	}

	refresh, err := jwtx.AssignRefresh(jwt.StandardClaims{
		// TODO: make `iss` and `sub` as env vars
		Issuer:    "v3nooom",
		IssuedAt:  time.Now().UTC().Unix(),
		Subject:   "st3llar",
		ExpiresAt: refreshExp.Unix(),
	})
	if err != nil {
		return nil, err
	}

	// sync token pairs
	token := &model.Token{
		Access:         access,
		Refresh:        refresh,
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
		return nil, errorx.Internal(err.Error())
	}

	// build response
	resp := dto.BuildRespCred(
		dto.WithAccount(req.Account),
		dto.WithOrganization(req.Organization),
		dto.WithEnvironment(req.Environment),
		dto.WithTokenPair(jwtx.Tokens{
			Token:   access,
			Refresh: refresh,
		}),
	)

	return resp, nil
}

func (cd *conductor) Refresh(c context.Context, req dto.ReqRefresh) (*dto.RespCredential, error) {
	token := new(model.Token)
	if err := db.Conn(c).
		Where(map[string]interface{}{
			"refresh": req.Refresh,
		}).
		First(token).Error; err != nil {
		if errors.As(err, &gorm.ErrRecordNotFound) {
			return nil, errorx.NotFound("none refresh token found")
		}

		return nil, errorx.Internal(err.Error())
	}

	jwtToken, _ := jwtx.ParseToken(token.Access)
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errorx.Internal("failed to assert JWT claims as MapClaims")
	}

	refreshExp := token.RefreshExpires.UTC()
	if refreshExp.Before(time.Now().UTC()) {
		return nil, errorx.UnauthorizedWithMsg("refresh token expired, please login again")
	}

	now := time.Now().UTC()
	tokenExp := now.AddDate(0, 1, 0)

	iss, err := jwtx.GetStrClaim(claims, "iss")
	if err != nil {
		return nil, err
	}
	sub, err := jwtx.GetStrClaim(claims, "sub")
	if err != nil {
		return nil, err
	}
	account, err := jwtx.GetStrClaim(claims, "account")
	if err != nil {
		return nil, err
	}
	organization, err := jwtx.GetStrClaim(claims, "organization")
	if err != nil {
		return nil, err
	}

	ac := jwtx.AccessClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    iss,
			Subject:   sub,
			IssuedAt:  now.Unix(),
			ExpiresAt: tokenExp.Unix(),
		},
		Account:      account,
		Organization: organization,
	}

	access, err := jwtx.AssignAccess(ac)
	if err != nil {
		return nil, err
	}

	// save tokens association
	token.Access = access
	token.AccessExpires = tokenExp
	token.UpdatedAt = &now

	if err := db.Conn(c).Save(token).Error; err != nil {
		return nil, errorx.Internal(err.Error())
	}

	resp := dto.BuildRespCred(
		dto.WithAccount(claims["account"].(string)),
		dto.WithOrganization(claims["organization"].(string)),
		dto.WithEnvironment(claims["environment"].(string)),
		dto.WithTokenPair(jwtx.Tokens{
			Token:   access,
			Refresh: req.Refresh,
		}),
	)

	return resp, nil
}

func (cd *conductor) Logout(c context.Context, req dto.ReqLogout) (*dto.RespLogout, error) {
	if err := db.Conn(c).
		Where(map[string]interface{}{
			"access": req.Token,
		}).
		Delete(&model.Token{}).Error; err != nil {
		return nil, errorx.Internal(err.Error())
	}

	return new(dto.RespLogout), nil
}
