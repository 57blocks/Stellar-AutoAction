package config

type (
	Credential struct {
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

	return nil, nil
}

func WriteCredential(path string, cred *Credential) error {

	return nil
}
