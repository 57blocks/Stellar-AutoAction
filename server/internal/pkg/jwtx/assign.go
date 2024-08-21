package jwtx

import (
	"encoding/base64"
	"time"

	"github.com/57blocks/auto-action/server/internal/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type (
	Claims struct {
		jwt.StandardClaims
		Account      string `json:"account"`
		Organization string `json:"organization"`
	}

	Tokens struct {
		_       struct{}
		Token   string `json:"token" toml:"token"`
		Refresh string `json:"refresh" toml:"refresh"`
	}

	AccessExpire  func() time.Time
	RefreshExpire func() time.Time
)

func Assign(ae, re time.Time) (*Tokens, error) {
	accessClaim := &Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "v3nooom",
			IssuedAt:  time.Now().UTC().Unix(),
			Subject:   "st3llar",
			ExpiresAt: ae.Unix(),
		},
		Account:      "Account_sample",
		Organization: "Organization_sample",
	}

	accessBytes, err := base64.StdEncoding.DecodeString(config.Global.JWT.PrivateKey)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	accessKey, err := jwt.ParseRSAPrivateKeyFromPEM(accessBytes)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	accessToken := jwt.NewWithClaims(jwt.GetSigningMethod(string(AlgRS256)), accessClaim)

	access, err := accessToken.SignedString(accessKey)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	refreshClaim := &jwt.StandardClaims{
		Issuer:    "v3nooom",
		IssuedAt:  time.Now().UTC().Unix(),
		Subject:   "st3llar",
		ExpiresAt: re.Unix(),
	}

	refreshBytes, err := base64.StdEncoding.DecodeString(config.Global.JWT.PrivateKey)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	refreshKey, err := jwt.ParseRSAPrivateKeyFromPEM(refreshBytes)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	refreshToken := jwt.NewWithClaims(jwt.GetSigningMethod(string(AlgRS256)), refreshClaim)

	refresh, err := refreshToken.SignedString(refreshKey)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return &Tokens{
		Token:   access,
		Refresh: refresh,
	}, nil
}
