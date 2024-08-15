package jwtx

import (
	"encoding/base64"
	"time"

	"github.com/57blocks/auto-action/server/internal/config"

	"github.com/dgrijalva/jwt-go"
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
)

func Assign() (*Tokens, error) {
	accessClaim := &Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "v3nooom",
			IssuedAt:  time.Now().UTC().Unix(),
			Subject:   "st3llar",
			ExpiresAt: time.Now().UTC().AddDate(0, 1, 0).Unix(),
		},
		Account:      "Account_sample",
		Organization: "Organization_sample",
	}

	accessBytes, err := base64.StdEncoding.DecodeString(config.Global.JWT.PrivateKey)
	if err != nil {
		return nil, err
	}

	accessKey, err := jwt.ParseRSAPrivateKeyFromPEM(accessBytes)
	if err != nil {
		return nil, err
	}

	accessToken := jwt.NewWithClaims(jwt.GetSigningMethod(string(AlgRS256)), accessClaim)

	access, err := accessToken.SignedString(accessKey)
	if err != nil {
		return nil, nil
	}

	refreshClaim := &jwt.StandardClaims{
		Issuer:    "v3nooom",
		IssuedAt:  time.Now().UTC().Unix(),
		Subject:   "st3llar",
		ExpiresAt: time.Now().UTC().AddDate(0, 3, 0).Unix(),
	}

	refreshBytes, err := base64.StdEncoding.DecodeString(config.Global.JWT.PrivateKey)
	if err != nil {
		return nil, err
	}

	refreshKey, err := jwt.ParseRSAPrivateKeyFromPEM(refreshBytes)
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.GetSigningMethod(string(AlgRS256)), refreshClaim)

	refresh, err := refreshToken.SignedString(refreshKey)
	if err != nil {
		return nil, nil
	}

	return &Tokens{
		Token:   access,
		Refresh: refresh,
	}, nil
}
