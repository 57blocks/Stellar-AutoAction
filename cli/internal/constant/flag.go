package constant

type FlagName string

// Flags for the configure command
const (
	FlagLog        FlagName = "log"
	FlagSource     FlagName = "source"
	FlagCredential FlagName = "credential"
	FlagEndPoint   FlagName = "endpoint"
)

// Flags for the signup command
const (
	FlagDescription FlagName = "description"
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

// FlagPayload Flags for Lambda invoke command
const (
	FlagPayload FlagName = "payload"
)

const (
	FlagFull FlagName = "full"
)

func (f FlagName) ValStr() string {
	return string(f)
}
