package oauth

type (
	Credential struct {
		Account      string `toml:"account" json:"account"`
		Organization string `toml:"organization" json:"organization"`
		*Bound       `toml:"bound" json:"bound"`
		*Tokens      `toml:"tokens" json:"tokens"`
	}
	CredOpt func(cred *Credential)

	Bound struct {
		_        struct{}
		Name     string `toml:"name" json:"name"`
		EndPoint string `toml:"endpoint" json:"endpoint"`
	}
	CredBoundOpt func(env *Bound)

	Tokens struct {
		_       struct{}
		Token   string `toml:"token" json:"token"`
		Refresh string `toml:"refresh" json:"refresh"`
	}
	CredTokenOpt func(tokens *Tokens)
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

func WithBoundEnv(env *Bound) CredOpt {
	return func(cred *Credential) {
		cred.Bound = env
	}
}

func WithTokens(session *Tokens) CredOpt {
	return func(cred *Credential) {
		cred.Tokens = session
	}
}

// BuildBound build the bound environment
func BuildBound(opts ...CredBoundOpt) *Bound {
	env := new(Bound)

	for _, opt := range opts {
		opt(env)
	}

	return env
}

func WithBoundName(name string) CredBoundOpt {
	return func(env *Bound) {
		env.Name = name
	}
}

func WithBoundEndPoint(endpoint string) CredBoundOpt {
	return func(env *Bound) {
		env.EndPoint = endpoint
	}
}

// BuildTokens build the credential's session pair
func BuildTokens(opts ...CredTokenOpt) *Tokens {
	session := new(Tokens)

	for _, opt := range opts {
		opt(session)
	}

	return session
}

func WithAccessToken(token string) CredTokenOpt {
	return func(tokens *Tokens) {
		tokens.Token = token
	}
}

func WithRefreshToken(refresh string) CredTokenOpt {
	return func(tokens *Tokens) {
		tokens.Refresh = refresh
	}
}
