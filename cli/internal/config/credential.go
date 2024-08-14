package config

type (
	Credential struct {
		Account      string `toml:"account" json:"account"`
		Organization string `toml:"organization" json:"organization"`
		*Environment `toml:"environment" json:"environment"`
		*Tokens      `toml:"tokens" json:"tokens"`
	}
	CredOpt func(cred *Credential)

	Environment struct {
		_        struct{}
		Name     string `toml:"name" json:"name"`
		EndPoint string `toml:"endpoint" json:"endpoint"`
	}
	CredEnvOpt func(env *Environment)

	Tokens struct {
		_       struct{}
		Token   string `toml:"token" json:"token"`
		Refresh string `toml:"refresh" json:"refresh"`
	}
	CredTokenOpt func(tokens *Tokens)
)

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

func WithEnvironment(env *Environment) CredOpt {
	return func(cred *Credential) {
		cred.Environment = env
	}
}

func WithSession(session *Tokens) CredOpt {
	return func(cred *Credential) {
		cred.Tokens = session
	}
}

// BuildCredEnv build the credential's environment pair
func BuildCredEnv(opts ...CredEnvOpt) *Environment {
	env := new(Environment)

	for _, opt := range opts {
		opt(env)
	}

	return env
}

func WithEnvName(name string) CredEnvOpt {
	return func(env *Environment) {
		env.Name = name
	}
}

func WithEnvEndPoint(endpoint string) CredEnvOpt {
	return func(env *Environment) {
		env.EndPoint = endpoint
	}
}

// BuildCredSession build the credential's session pair
func BuildCredSession(opts ...CredTokenOpt) *Tokens {
	session := new(Tokens)

	for _, opt := range opts {
		opt(session)
	}

	return session
}

func WithSessionToken(token string) CredTokenOpt {
	return func(tokens *Tokens) {
		tokens.Token = token
	}
}

func WithSessionRefreshToken(refresh string) CredTokenOpt {
	return func(tokens *Tokens) {
		tokens.Refresh = refresh
	}
}
