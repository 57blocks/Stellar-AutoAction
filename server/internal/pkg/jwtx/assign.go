package jwtx

import (
	"encoding/base64"
	"fmt"

	"github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	pkgLog "github.com/57blocks/auto-action/server/internal/pkg/log"

	"github.com/dgrijalva/jwt-go"
)

type (
	AccessClaims struct {
		jwt.StandardClaims
		Account      string `json:"account"`
		Organization string `json:"organization"`
		Environment  string `json:"environment"`
	}

	Tokens struct {
		_       struct{}
		Token   string `json:"token" toml:"token"`
		Refresh string `json:"refresh" toml:"refresh"`
	}

	ClaimPair struct {
		Token   AccessClaims
		Refresh jwt.StandardClaims
	}
)

func AssignAccess(accClaim AccessClaims) (string, error) {
	priPEM, err := base64.StdEncoding.DecodeString(config.GlobalConfig.JWT.PrivateKey)
	if err != nil {
		return "", errorx.Internal("decode private key failed")
	}

	priKey, err := jwt.ParseRSAPrivateKeyFromPEM(priPEM)
	if err != nil {
		return "", errorx.Internal("parse private key failed")
	}

	accToken := jwt.NewWithClaims(jwt.GetSigningMethod(string(AlgRS256)), accClaim)

	access, err := accToken.SignedString(priKey)
	if err != nil {
		return "", errorx.Internal("sign access token failed")
	}

	return access, nil
}

func AssignRefresh(refClaim jwt.StandardClaims) (string, error) {
	priPEM, err := base64.StdEncoding.DecodeString(config.GlobalConfig.JWT.PrivateKey)
	if err != nil {
		return "", errorx.Internal("decode private key failed")
	}

	priKey, err := jwt.ParseRSAPrivateKeyFromPEM(priPEM)
	if err != nil {
		return "", errorx.Internal("parse private key failed")
	}

	refToken := jwt.NewWithClaims(jwt.GetSigningMethod(string(AlgRS256)), refClaim)

	refresh, err := refToken.SignedString(priKey)
	if err != nil {
		return "", errorx.Internal("sign access token failed")
	}

	return refresh, nil
}

func ParseToken(tokenStr string) (*jwt.Token, error) {
	pubPEM, err := base64.StdEncoding.DecodeString(config.GlobalConfig.JWT.PublicKey)
	if err != nil {
		return nil, errorx.Internal("decode private key failed")
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubPEM)
	if err != nil {
		return nil, errorx.Internal("parse public key failed")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return pubKey, nil
	})
	if err != nil {
		return nil, errorx.Internal("parse token failed")
	}

	if !token.Valid || token.Claims.Valid() != nil {
		return nil, errorx.Unauthorized()
	}

	return token, nil
}

// GetStrClaim extracts a string claim from jwt.MapClaims
func GetStrClaim(claims jwt.MapClaims, key string) (string, error) {
	value, ok := claims[key]
	if !ok {
		pkgLog.Logger.ERROR("claim not found", map[string]interface{}{"claim_key": key})
		return "", errorx.Internal(fmt.Sprintf("claim not found by: %v\n", key))
	}

	strValue, ok := value.(string)
	if !ok {
		pkgLog.Logger.ERROR(
			"claim value conversion error",
			map[string]interface{}{"claim_key": key, "value": value},
		)
		return "", errorx.Internal(fmt.Sprintf("claim value: %v, conversion error", value))
	}

	return strValue, nil
}
