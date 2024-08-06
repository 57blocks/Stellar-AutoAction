package config

type (
	Credential struct {
		Account      string `toml:"account"`
		Organization string `toml:"organization"`
		*Environment `toml:"environment"`
		*Session     `toml:"session"`
	}
	CredOpt func(cred *Credential)

	Environment struct {
		_        struct{}
		Name     string `toml:"name"`
		EndPoint string `toml:"endpoint"`
	}
	CredEnvOpt func(env *Environment)

	Session struct {
		_            struct{}
		Token        string `toml:"token"`
		RefreshToken string `toml:"refresh_token"`
	}
	CredSessionOpt func(session *Session)
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

func WithSession(session *Session) CredOpt {
	return func(cred *Credential) {
		cred.Session = session
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
func BuildCredSession(opts ...CredSessionOpt) *Session {
	session := new(Session)

	for _, opt := range opts {
		opt(session)
	}

	return session
}

func WithSessionToken(token string) CredSessionOpt {
	return func(session *Session) {
		session.Token = token
	}
}

func WithSessionRefreshToken(refreshToken string) CredSessionOpt {
	return func(session *Session) {
		session.RefreshToken = refreshToken
	}
}
