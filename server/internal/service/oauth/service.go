package oauth

import (
	"context"
	"time"

	"github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/dto"
	"github.com/57blocks/auto-action/server/internal/model"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/57blocks/auto-action/server/internal/repo"
	"github.com/57blocks/auto-action/server/internal/third-party/decrypt"
	"github.com/57blocks/auto-action/server/internal/third-party/jwtx"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type (
	Service interface {
		Login(c context.Context, req dto.ReqLogin) (*dto.RespCredential, error)
		Refresh(c context.Context, raw string) (*dto.RespCredential, error)
		Logout(c context.Context, raw string) (*dto.RespLogout, error)
	}
	service struct {
		jwtx      jwtx.JWT
		decrypter decrypt.Decrypter
		oauthRepo repo.OAuth
	}
)

func NewOAuthService() {
	if ServiceImpl == nil {
		repo.NewOAuth()

		ServiceImpl = &service{
			jwtx:      jwtx.RS256,
			decrypter: decrypt.RSADecrypter,
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

	rawPwdBytes, err := svc.decrypter.Decrypt([]byte(req.Password))
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), rawPwdBytes); err != nil {
		return nil, errorx.BadRequest("password not match")
	}

	// tokens assignment
	now := time.Now().UTC()
	accessID := svc.jwtx.GenerateID()
	accessExp := now.AddDate(0, 0, 7)
	refreshID := svc.jwtx.GenerateID()
	refreshExp := now.AddDate(0, 1, 0)

	access, err := svc.jwtx.Assign(&jwtx.AAClaims{
		StdJWTClaims: jwt.StandardClaims{
			Audience:  config.GlobalConfig.EndPoint,
			ExpiresAt: accessExp.Unix(),
			Id:        accessID,
			IssuedAt:  now.Unix(),
			Issuer:    req.Organization,
			NotBefore: now.Unix(),
			Subject:   u.Account,
		},
	})
	if err != nil {
		return nil, err
	}

	refresh, err := svc.jwtx.Assign(&jwtx.AAClaims{
		StdJWTClaims: jwt.StandardClaims{
			Audience:  config.GlobalConfig.EndPoint,
			ExpiresAt: refreshExp.Unix(),
			Id:        refreshID,
			IssuedAt:  now.Unix(),
			Issuer:    req.Organization,
			NotBefore: now.Unix(),
			//NotBefore: accessExp.Unix(), // won't be valid until access token expires
			Subject: u.Account,
		},
	})
	if err != nil {
		return nil, err
	}

	// sync token pairs
	token := &model.Token{
		UserId:         u.ID,
		Access:         access,
		AccessID:       accessID,
		AccessExpires:  accessExp,
		Refresh:        refresh,
		RefreshID:      refreshID,
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
		dto.WithTokenPair(jwtx.TokenPair{
			Access:  access,
			Refresh: refresh,
		}),
	)

	return resp, nil
}

func (svc *service) Refresh(c context.Context, raw string) (*dto.RespCredential, error) {
	jwtClaims, err := svc.jwtx.Parse(raw)
	if err != nil {
		return nil, err
	}

	aaClaims, ok := jwtClaims.(*jwtx.AAClaims)
	if !ok {
		return nil, errorx.Internal("failed to assert JWT claims as AAClaims")
	}

	token, err := svc.oauthRepo.FindTokenByRefreshID(c, aaClaims.StdJWTClaims.Id) // use refresh id
	if err != nil {
		return nil, errorx.UnauthorizedWithMsg("invalid refresh token")
	}

	now := time.Now().UTC()
	accessID := svc.jwtx.GenerateID()
	accessExp := now.AddDate(0, 0, 7)

	access, err := svc.jwtx.Assign(&jwtx.AAClaims{
		StdJWTClaims: jwt.StandardClaims{
			Audience:  config.GlobalConfig.EndPoint,
			ExpiresAt: accessExp.Unix(),
			Id:        accessID,
			IssuedAt:  now.Unix(),
			Issuer:    aaClaims.StdJWTClaims.Issuer,
			NotBefore: now.Unix(),
			Subject:   aaClaims.StdJWTClaims.Subject,
		},
	})
	if err != nil {
		return nil, err
	}

	// save tokens association
	token.Access = access
	token.AccessID = accessID
	token.AccessExpires = accessExp
	token.UpdatedAt = &now

	if err := svc.oauthRepo.SyncToken(c, token); err != nil {
		return nil, err
	}

	resp := dto.BuildRespCred(
		dto.WithAccount(aaClaims.StdJWTClaims.Subject),
		dto.WithOrganization(aaClaims.StdJWTClaims.Issuer),
		dto.WithEnvironment(config.GlobalConfig.Name),
		dto.WithTokenPair(jwtx.TokenPair{
			Access:  access,
			Refresh: raw,
		}),
	)

	return resp, nil
}

func (svc *service) Logout(c context.Context, raw string) (*dto.RespLogout, error) {
	if err := svc.oauthRepo.DeleteTokenByAccess(c, raw); err != nil {
		return nil, err
	}

	return new(dto.RespLogout), nil
}
