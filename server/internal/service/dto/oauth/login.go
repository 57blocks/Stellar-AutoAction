package oauth

import (
	"github.com/57blocks/auto-action/server/internal/pkg/jwtx"
)

type (
	ReqLogin struct {
		_            struct{}
		Account      string `json:"account"`
		Organization string `json:"organization"`
		Password     []byte `json:"password"`
		Environment  string `json:"environment"`
	}

	RespLogin struct {
		_            struct{}
		Account      string `json:"account" toml:"account"`
		Organization string `json:"organization" toml:"organization"`
		Environment  string `toml:"environment" json:"environment"`
		*jwtx.Tokens `json:"tokens" toml:"tokens"`
	}
	RespOpt func(cred *RespLogin)
)

// BuildRespLogin build the RespLogin of Login
func BuildRespLogin(opts ...RespOpt) *RespLogin {
	cred := new(RespLogin)

	for _, opt := range opts {
		opt(cred)
	}

	return cred
}

func WithAccount(account string) RespOpt {
	return func(resp *RespLogin) {
		resp.Account = account
	}
}

func WithOrganization(organization string) RespOpt {
	return func(resp *RespLogin) {
		resp.Organization = organization
	}
}

func WithEnvironment(environment string) RespOpt {
	return func(resp *RespLogin) {
		resp.Environment = environment
	}
}

func WithTokenPair(tokens *jwtx.Tokens) RespOpt {
	return func(resp *RespLogin) {
		resp.Tokens = tokens
	}
}
