package constant

type FlagName string

// Flags for the configure command
const (
	FlagPrefix     FlagName = "prefix"
	FlagLog        FlagName = "log"
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
	FlagCron FlagName = "cron"
	FlagRate FlagName = "rate"
)

// FlagPayload Flags for Lambda invoke command
const (
	FlagPayload FlagName = "payload"
)

func (f FlagName) ValStr() string {
	return string(f)
}
