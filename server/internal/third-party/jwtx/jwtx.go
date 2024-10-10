package jwtx

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"

	"github.com/dgrijalva/jwt-go"
)

//go:generate mockgen -destination ../../testdata/jwtx_mock.go -package testdata -source jwtx.go JWT
type (
	JWT interface {
		Assigner
		Parser
		Generator
	}
	Parser interface {
		Parse(raw string) (jwt.Claims, error)
	}
	Assigner interface {
		Assign(claim jwt.Claims) (string, error)
	}
	Generator interface {
		GenerateID() string
	}
)

var RS256 JWT

func Setup() error {
	if RS256 == nil {
		RS256 = &rs256{}
	}

	return nil
}

type (
	AAClaims struct {
		_ struct{}
		// Aud - server host, which could distinguish different environments
		// Iss - organization name
		// Sub - account name
		StdJWTClaims jwt.StandardClaims
	}

	TokenPair struct {
		Access  string
		Refresh string
	}
)

func (aac *AAClaims) Valid() error {
	if err := aac.StdJWTClaims.Valid(); err != nil {
		return err
	}

	// TODO: confirm with that using Name or EndPoint when the terraform is ready
	if aac.StdJWTClaims.Audience != config.GlobalConfig.Bound.EndPoint {
		return errorx.UnauthorizedWithMsg("invalid audience")
	}

	return nil
}

type rs256 struct{}

func (rs *rs256) Assign(claim jwt.Claims) (string, error) {
	pem, err := base64.StdEncoding.DecodeString(config.GlobalConfig.JWT.PrivateKey)
	if err != nil {
		return "", errorx.Internal("decode private key failed")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(pem)
	if err != nil {
		return "", errorx.Internal("parse private key failed")
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(config.GlobalConfig.JWT.Protocol), claim)

	signed, err := token.SignedString(privateKey)
	if err != nil {
		return "", errorx.Internal("sign access token failed")
	}

	return signed, nil
}

func (rs *rs256) Parse(raw string) (jwt.Claims, error) {
	pem, err := base64.StdEncoding.DecodeString(config.GlobalConfig.JWT.PublicKey)
	if err != nil {
		return nil, errorx.Internal("decode public key failed")
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pem)
	if err != nil {
		return nil, errorx.Internal("parse public key failed")
	}

	claims := new(AAClaims)
	token, err := jwt.ParseWithClaims(raw, claims, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		return nil, errorx.Internal(fmt.Sprintf("parse token failed: %s", err.Error()))
	}
	if !token.Valid {
		return nil, errorx.Unauthorized()
	}
	if err := token.Claims.Valid(); err != nil {
		return nil, err
	}

	return claims, nil
}

const length = 16

func (rs *rs256) GenerateID() string {
	randomBytes := make([]byte, length)
	rand.Read(randomBytes)
	randomString := base64.URLEncoding.EncodeToString(randomBytes)

	return randomString[:length]
}
