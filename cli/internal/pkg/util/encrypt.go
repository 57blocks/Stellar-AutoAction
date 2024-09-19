package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"

	"github.com/57blocks/auto-action/cli/internal/constant"
	"github.com/57blocks/auto-action/cli/internal/pkg/errorx"
)

func LoadPublicKey() (*rsa.PublicKey, error) {
	pubPEMBytes, err := base64.StdEncoding.DecodeString(constant.PublicKey)
	if err != nil {
		return nil, errorx.Internal("decode public key failed")
	}

	block, _ := pem.Decode(pubPEMBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errorx.Internal("failed to decode PEM block containing public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pubKey := pubKey.(type) {
	case *rsa.PublicKey:
		return pubKey, nil
	default:
		return nil, errorx.Internal("not an RSA public key")
	}
}

func EncryptPassword(password string, publicKey *rsa.PublicKey) (string, error) {
	encryptedBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, []byte(password), nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}
