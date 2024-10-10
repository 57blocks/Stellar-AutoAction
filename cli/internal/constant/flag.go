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
	FlagOrganization FlagName = "organization"
)

// Flags for Action register command
const (
	FlagAt   FlagName = "at"
	FlagCron FlagName = "cron"
	FlagRate FlagName = "rate"
)

// FlagPayload Flags for Action invoke command
const (
	FlagPayload FlagName = "payload"
)

const (
	FlagFull FlagName = "full"
)

func (f FlagName) ValStr() string {
	return string(f)
}
