package dto

import (
	"github.com/57blocks/auto-action/server/internal/third-party/jwtx"
)

// Login related dto
type (
	ReqLogin struct {
		_            struct{}
		Account      string `json:"account"`
		Organization string `json:"organization"`
		Password     []byte `json:"password"`
		Environment  string `json:"environment"`
	}

	RespCredential struct {
		_            struct{}
		Account      string `json:"account" toml:"account"`
		Organization string `json:"organization" toml:"organization"`
		Environment  string `toml:"environment" json:"environment"`
		jwtx.Tokens  `json:"tokens" toml:"tokens"`
	}
	RespCredOpt func(cred *RespCredential)
)

// BuildRespCred build the RespCredential of Login
func BuildRespCred(opts ...RespCredOpt) *RespCredential {
	cred := new(RespCredential)

	for _, opt := range opts {
		opt(cred)
	}

	return cred
}

func WithAccount(account string) RespCredOpt {
	return func(resp *RespCredential) {
		resp.Account = account
	}
}

func WithOrganization(organization string) RespCredOpt {
	return func(resp *RespCredential) {
		resp.Organization = organization
	}
}

func WithEnvironment(environment string) RespCredOpt {
	return func(resp *RespCredential) {
		resp.Environment = environment
	}
}

func WithTokenPair(tokens jwtx.Tokens) RespCredOpt {
	return func(resp *RespCredential) {
		resp.Tokens = tokens
	}
}

// Logout related dto
type (
	ReqLogout struct {
		_     struct{}
		Token string `json:"token"`
	}

	RespLogout struct{}
)

// ReqRefresh related dto
type ReqRefresh struct {
	_       struct{}
	Refresh string `json:"refresh"`
}
