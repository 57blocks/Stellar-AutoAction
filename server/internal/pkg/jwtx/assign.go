package jwtx

import (
	"encoding/base64"
	"github.com/57blocks/auto-action/server/internal/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type (
	TokenClaims struct {
		jwt.StandardClaims
		Account      string `json:"account"`
		Organization string `json:"organization"`
	}

	Tokens struct {
		_       struct{}
		Token   string `json:"token" toml:"token"`
		Refresh string `json:"refresh" toml:"refresh"`
	}

	ClaimPair struct {
		Token   TokenClaims
		Refresh jwt.StandardClaims
	}
)

func Assign(cPair ClaimPair) (*Tokens, error) {
	accessBytes, err := base64.StdEncoding.DecodeString(config.Global.JWT.PrivateKey)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	accessKey, err := jwt.ParseRSAPrivateKeyFromPEM(accessBytes)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	accessToken := jwt.NewWithClaims(jwt.GetSigningMethod(string(AlgRS256)), cPair.Token)

	access, err := accessToken.SignedString(accessKey)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	refreshBytes, err := base64.StdEncoding.DecodeString(config.Global.JWT.PrivateKey)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	refreshKey, err := jwt.ParseRSAPrivateKeyFromPEM(refreshBytes)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	refreshToken := jwt.NewWithClaims(jwt.GetSigningMethod(string(AlgRS256)), cPair.Refresh)

	refresh, err := refreshToken.SignedString(refreshKey)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return &Tokens{
		Token:   access,
		Refresh: refresh,
	}, nil
}

func Parse(token string) (*jwt.Token, error) {
	accessBytes, err := base64.StdEncoding.DecodeString(config.Global.JWT.PublicKey)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	accessKey, err := jwt.ParseRSAPublicKeyFromPEM(accessBytes)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return accessKey, nil
	})
}
