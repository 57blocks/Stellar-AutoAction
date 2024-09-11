package constant

type FlagName string

// Flags for the configure command
const (
	FlagLog        FlagName = "logx"
	FlagSource     FlagName = "source"
	FlagCredential FlagName = "credential"
	FlagEndPoint   FlagName = "endpoint"
)

// Flags for the login command
const (
	FlagAccount      FlagName = "account"
	FlagEnvironment  FlagName = "environment"
	FlagOrganization FlagName = "organization"
)

// Flags for Lambda register command
const (
	FlagAt   FlagName = "at"
	FlagCron FlagName = "cron"
	FlagRate FlagName = "rate"
)

// FlagEnv Flags for Wallet verify command
const (
	FlagEnv FlagName = "env"
)

func (f FlagName) ValStr() string {
	return string(f)
}
