package oauth

import (
	"context"
	"time"

	"github.com/57blocks/auto-action/server/internal/db"
	"github.com/57blocks/auto-action/server/internal/model"
	"github.com/57blocks/auto-action/server/internal/pkg/jwtx"
	oauth2 "github.com/57blocks/auto-action/server/internal/service/dto/oauth"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	Service interface {
		Login(c context.Context, req oauth2.ReqLogin) (*oauth2.RespCredential, error)
		Refresh(c context.Context, req oauth2.ReqRefresh) (*oauth2.RespCredential, error)
		Logout(c context.Context, req oauth2.ReqLogout) (*oauth2.RespLogout, error)
	}
	conductor struct{}
)

var Conductor Service

func init() {
	if Conductor == nil {
		Conductor = &conductor{}
	}
}

func (cd *conductor) Login(c context.Context, req oauth2.ReqLogin) (*oauth2.RespCredential, error) {
	user := new(model.User)
	if err := db.Conn(c).Table(user.TableNameWithAbbr()).
		Joins("LEFT JOIN organization AS o ON u.organization_id = o.id").
		Where(map[string]interface{}{
			"u.account": req.Account,
			"o.name":    req.Organization,
		}).
		First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrap(err, "user/organization not found")
		}
		return nil, errors.Wrap(err, err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), req.Password); err != nil {
		return nil, errors.New("password not match")
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
		return nil, errors.New(err.Error())
	}

	// build response
	resp := oauth2.BuildRespCred(
		oauth2.WithAccount(req.Account),
		oauth2.WithOrganization(req.Organization),
		oauth2.WithEnvironment(req.Environment),
		oauth2.WithTokenPair(jwtx.Tokens{
			Token:   access,
			Refresh: refresh,
		}),
	)

	return resp, nil
}

func (cd *conductor) Refresh(c context.Context, req oauth2.ReqRefresh) (*oauth2.RespCredential, error) {
	token := new(model.Token)
	if err := db.Conn(c).
		Where(map[string]interface{}{
			"refresh": req.Refresh,
		}).
		First(token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("refresh token not found")
		}
		return nil, errors.New(err.Error())
	}

	jwtToken, _ := jwtx.ParseToken(token.Access)
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to assert JWT claims as MapClaims")
	}

	refreshExp := token.RefreshExpires.UTC()
	if refreshExp.Before(time.Now().UTC()) {
		return nil, errors.New("refresh token expired, please login again")
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
		return nil, errors.New(err.Error())
	}

	resp := oauth2.BuildRespCred(
		oauth2.WithAccount(claims["account"].(string)),
		oauth2.WithOrganization(claims["organization"].(string)),
		oauth2.WithEnvironment(claims["environment"].(string)),
		oauth2.WithTokenPair(jwtx.Tokens{
			Token:   access,
			Refresh: req.Refresh,
		}),
	)

	return resp, nil
}

func (cd *conductor) Logout(c context.Context, req oauth2.ReqLogout) (*oauth2.RespLogout, error) {
	if err := db.Conn(c).
		Where(map[string]interface{}{
			"access": req.Token,
		}).
		Delete(&model.Token{}).Error; err != nil {
		return nil, errors.New(err.Error())
	}

	return new(oauth2.RespLogout), nil
}
