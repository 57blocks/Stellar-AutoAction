package oauth

import (
	"context"
	"time"

	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/model"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/57blocks/auto-action/server/internal/repo"
	"github.com/57blocks/auto-action/server/internal/third-party/jwtx"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type (
	Service interface {
		Login(c context.Context, req dto.ReqLogin) (*dto.RespCredential, error)
		Refresh(c context.Context, req dto.ReqRefresh) (*dto.RespCredential, error)
		Logout(c context.Context, req dto.ReqLogout) (*dto.RespLogout, error)
	}
	service struct {
		oauthRepo repo.OAuth
	}
)

func NewOAuthService() {
	if ServiceImpl == nil {
		repo.NewOAuth()

		ServiceImpl = &service{
			oauthRepo: repo.OAuthRepo,
		}
	}
}

func (svc *service) Login(c context.Context, req dto.ReqLogin) (*dto.RespCredential, error) {
	u, err := svc.oauthRepo.FindUserByOrgAcn(c, &dto.ReqOrgAcn{
		OrgName: req.Organization,
		AcnName: req.Account,
	})
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), req.Password); err != nil {
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
		UserId:         u.ID,
		AccessExpires:  tokenExp,
		RefreshExpires: refreshExp,
	}
	if err := svc.oauthRepo.SyncToken(c, token); err != nil {
		return nil, err
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

func (svc *service) Refresh(c context.Context, req dto.ReqRefresh) (*dto.RespCredential, error) {
	token, err := svc.oauthRepo.FindTokenByRefresh(c, req.Refresh)
	if err != nil {
		return nil, errorx.UnauthorizedWithMsg("invalid refresh token")
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

	if err := svc.oauthRepo.SyncToken(c, token); err != nil {
		return nil, err
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

func (svc *service) Logout(c context.Context, req dto.ReqLogout) (*dto.RespLogout, error) {
	if err := svc.oauthRepo.DeleteTokenByAccess(c, req.Token); err != nil {
		return nil, err
	}

	return new(dto.RespLogout), nil
}
