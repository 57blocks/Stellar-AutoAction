package jwtx

import (
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
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "v3nooom",
			IssuedAt:  time.Now().UTC().Unix(),
			Subject:   "st3llar",
			ExpiresAt: time.Now().UTC().AddDate(0, 1, 0).Unix(),
		},
		Account:      "Account_sample",
		Organization: "Organization_sample",
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(config.Global.JWT.PrivateKey))
	if err != nil {
		return nil, err
	}

	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod(string(AlgRS256)), claims)

	token, err := jwtToken.SignedString(signKey)
	if err != nil {
		return nil, nil
	}

	return &Tokens{
		Token:   token,
		Refresh: "Refresh_sample",
	}, nil
}
