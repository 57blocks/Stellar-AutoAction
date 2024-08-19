package oauth

import "github.com/57blocks/auto-action/server/internal/pkg/jwtx"

// BuildResp build the RespLogin
func BuildResp(opts ...RespLoginOpt) *RespLogin {
	cred := new(RespLogin)

	for _, opt := range opts {
		opt(cred)
	}

	return cred
}

func WithAccount(account string) RespLoginOpt {
	return func(resp *RespLogin) {
		resp.Account = account
	}
}

func WithOrganization(organization string) RespLoginOpt {
	return func(resp *RespLogin) {
		resp.Organization = organization
	}
}

func WithBound(bound *Bound) RespLoginOpt {
	return func(resp *RespLogin) {
		resp.Bound = bound
	}
}

func WithTokens(tokens *jwtx.Tokens) RespLoginOpt {
	return func(resp *RespLogin) {
		resp.Tokens = tokens
	}
}

// BuildBound build the bound environment
func BuildBound(opts ...RespBoundOpt) *Bound {
	bound := new(Bound)

	for _, opt := range opts {
		opt(bound)
	}

	return bound
}

func WithBoundName(name string) RespBoundOpt {
	return func(env *Bound) {
		env.Name = name
	}
}

func WithBoundEndPoint(endpoint string) RespBoundOpt {
	return func(env *Bound) {
		env.EndPoint = endpoint
	}
}
