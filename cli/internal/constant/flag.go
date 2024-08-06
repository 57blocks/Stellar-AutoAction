package constant

type FlagName string

const (
	FlagAccount      FlagName = "account"
	FlagCredential   FlagName = "credential"
	FlagEnvPrefix    FlagName = "env-prefix"
	FlagEnvironment  FlagName = "environment"
	FlagLogLevel     FlagName = "log-level"
	FlagOrganization FlagName = "organization"
)

func (f FlagName) ValStr() string {
	return string(f)
}
