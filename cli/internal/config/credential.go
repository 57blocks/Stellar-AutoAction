package config

import (
	"fmt"
	"os"

	"github.com/57blocks/auto-action/cli/internal/pkg/errorx"

	"github.com/BurntSushi/toml"
)

type (
	Credential struct {
		_            struct{}
		Account      string `toml:"account" json:"account"`
		Organization string `toml:"organization" json:"organization"`
		*Tokens      `toml:"tokens" json:"tokens"`
	}
	CredOpt func(cred *Credential)

	Tokens struct {
		_       struct{}
		Access  string `toml:"access" json:"access"`
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

func WithAccess(access string) CredOpt {
	return func(cred *Credential) {
		cred.Access = access
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
		return nil, errorx.Internal(fmt.Sprintf("read credential error: %s", err.Error()))
	}

	cred := new(Credential)

	if _, err := toml.Decode(string(data), cred); err != nil {
		return nil, errorx.Internal(fmt.Sprintf("decode credential error: %s", err.Error()))
	}

	return cred, nil
}

func WriteCredential(path string, cred *Credential) error {
	tomlBytes, err := toml.Marshal(cred)
	if err != nil {
		return errorx.Internal(fmt.Sprintf("marshal credential error: %s", err))
	}

	if err := os.WriteFile(path, tomlBytes, 0666); err != nil {
		return errorx.Internal(fmt.Sprintf("write credential error: %s", err.Error()))
	}

	return nil
}

func RemoveCredential(path string) error {
	if err := os.Remove(path); err != nil {
		return errorx.Internal(fmt.Sprintf("remove credential error: %s", err.Error()))
	}

	return nil
}

func Token() (string, error) {
	cfg, _ := ReadConfig()

	credential, err := ReadCredential(cfg.Credential)
	if err != nil {
		return "", err
	}

	return credential.Access, nil
}
