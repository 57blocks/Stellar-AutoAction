package decrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"

	configx "github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
)

//go:generate mockgen -destination ./decrypt_mock.go -package decrypt -source decrypt.go Decrypter
type Decrypter interface {
	Decrypt(data []byte) ([]byte, error)
}

type rsaDecrypter struct{}

var RSADecrypter Decrypter

func Setup() error {
	if RSADecrypter == nil {
		RSADecrypter = &rsaDecrypter{}
	}

	return nil
}

func (d *rsaDecrypter) Decrypt(data []byte) ([]byte, error) {
	privateKey, err := loadPrivateKey()
	if err != nil {
		return nil, err
	}

	return decryptPassword(string(data), privateKey)
}

func loadPrivateKey() (*rsa.PrivateKey, error) {
	priPEMBytes, err := base64.StdEncoding.DecodeString(configx.GlobalConfig.RSA.PrivateKey)
	if err != nil {
		return nil, errorx.Internal("decode public key failed")
	}

	block, _ := pem.Decode(priPEMBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errorx.Internal("failed to decode RSA block containing private key")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func decryptPassword(encryptedPassword string, privateKey *rsa.PrivateKey) ([]byte, error) {
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedPassword)
	if err != nil {
		return nil, err
	}

	decryptedBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encryptedBytes, nil)
	if err != nil {
		return nil, err
	}

	return decryptedBytes, nil
}
