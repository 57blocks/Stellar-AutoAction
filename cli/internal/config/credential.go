package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"os"
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
		return nil, err
	}

	cred := new(Credential)

	if _, err := toml.Decode(string(data), cred); err != nil {
		return nil, errors.New(fmt.Sprintf("reading credential error: %s\n", err.Error()))
	}

	return cred, nil
}

func WriteCredential(path string, cred *Credential) error {
	tomlBytes, err := toml.Marshal(cred)
	if err != nil {
		return fmt.Errorf("marshalling credential error: %w", err)
	}

	if err := os.WriteFile(path, tomlBytes, 0666); err != nil {
		return errors.New(fmt.Sprintf("writing credential error: %s\n", err.Error()))
	}

	return nil
}

func RemoveCredential(path string) error {
	if err := os.Remove(path); err != nil {
		return errors.New(fmt.Sprintf("removing credential error: %s\n", err.Error()))
	}

	return nil
}
