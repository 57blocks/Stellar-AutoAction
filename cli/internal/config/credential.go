package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type (
	Credential struct {
		_            struct{}
		Account      string `toml:"account" json:"account"`
		Organization string `toml:"organization" json:"organization"`
		Environment  string `toml:"environment" json:"environment"`
		*Tokens      `toml:"tokens" json:"tokens"`
	}
	CredOpt func(cred *Credential)

	Tokens struct {
		_       struct{}
		Token   string `toml:"token" json:"token"`
		Refresh string `toml:"refresh" json:"refresh"`
	}
)

// BuildCred build the credential pair
func BuildCred(opts ...CredOpt) *Credential {
	cred := new(Credential)

	for _, opt := range opts {
		opt(cred)
	}

	return cred
}

func WithAccount(account string) CredOpt {
	return func(cred *Credential) {
		cred.Account = account
	}
}

func WithOrganization(organization string) CredOpt {
	return func(cred *Credential) {
		cred.Organization = organization
	}
}

func WithEnvironment(env string) CredOpt {
	return func(cred *Credential) {
		cred.Environment = env
	}
}

func WithAccess(access string) CredOpt {
	return func(cred *Credential) {
		cred.Token = access
	}
}

func WithRefresh(refresh string) CredOpt {
	return func(cred *Credential) {
		cred.Refresh = refresh
	}
}

func ReadCredential(path string) (*Credential, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		errMsg := fmt.Sprintf("read credential error: %s\n", err.Error())
		slog.Error(errMsg)
		return nil, errors.Wrap(err, errMsg)
	}

	cred := new(Credential)

	if _, err := toml.Decode(string(data), cred); err != nil {
		errMsg := fmt.Sprintf("decode credential error: %s\n", err.Error())
		slog.Error(errMsg)
		return nil, errors.Wrap(err, errMsg)
	}

	return cred, nil
}

func WriteCredential(path string, cred *Credential) error {
	tomlBytes, err := toml.Marshal(cred)
	if err != nil {
		errMsg := fmt.Sprintf("marshal credential error: %s\n", err)
		slog.Error(errMsg)
		return errors.Wrap(err, errMsg)
	}

	if err := os.WriteFile(path, tomlBytes, 0666); err != nil {
		errMsg := fmt.Sprintf("write credential error: %s\n", err.Error())
		slog.Error(errMsg)
		return errors.Wrap(err, errMsg)
	}

	return nil
}

func RemoveCredential(path string) error {
	if err := os.Remove(path); err != nil {
		errMsg := fmt.Sprintf("remove credential error: %s\n", err.Error())
		slog.Error(errMsg)
		return errors.Wrap(err, errMsg)
	}

	return nil
}

func Token() (string, error) {
	cfg, _ := ReadConfig()

	credential, err := ReadCredential(cfg.Credential)
	if err != nil {
		return "", err
	}

	return credential.Token, nil
}
