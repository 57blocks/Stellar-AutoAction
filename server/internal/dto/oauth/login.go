package oauth

import "github.com/57blocks/auto-action/server/internal/pkg/jwtx"

type (
	Request struct {
		_            struct{}
		Account      string `json:"account"`
		Organization string `json:"organization"`
		Password     []byte `json:"password"`
		Environment  string `json:"environment"`
	}

	Response struct {
		_            struct{}
		Account      string `json:"account" toml:"account"`
		Organization string `json:"organization" toml:"organization"`
		*jwtx.Tokens `json:"tokens" toml:"tokens"`
		*Bound       `json:"bound" toml:"bound"`
	}
	ResponseOpt func(cred *Response)

	Bound struct {
		_        struct{}
		Name     string `json:"name" toml:"name"`
		EndPoint string `json:"endpoint" toml:"endpoint"`
	}
	BoundOpt func(bound *Bound)
)

// BuildResp build the Response
func BuildResp(opts ...ResponseOpt) *Response {
	cred := new(Response)

	for _, opt := range opts {
		opt(cred)
	}

	return cred
}

func WithAccount(account string) ResponseOpt {
	return func(resp *Response) {
		resp.Account = account
	}
}

func WithOrganization(organization string) ResponseOpt {
	return func(resp *Response) {
		resp.Organization = organization
	}
}

func WithBound(bound *Bound) ResponseOpt {
	return func(resp *Response) {
		resp.Bound = bound
	}
}

func WithTokens(tokens *jwtx.Tokens) ResponseOpt {
	return func(resp *Response) {
		resp.Tokens = tokens
	}
}

// BuildBound build the bound environment
func BuildBound(opts ...BoundOpt) *Bound {
	bound := new(Bound)

	for _, opt := range opts {
		opt(bound)
	}

	return bound
}

func WithBoundName(name string) BoundOpt {
	return func(env *Bound) {
		env.Name = name
	}
}

func WithBoundEndPoint(endpoint string) BoundOpt {
	return func(env *Bound) {
		env.EndPoint = endpoint
	}
}
