package constant

// Config is the type of strings in config
type Config string

const (
	ConfigurationType Config = "toml"
	ConfigName        Config = ".autoaction"
	CredentialName    Config = ".autoaction-credential"
)

func (cc Config) ValStr() string {
	return string(cc)
}

// Log related types in config
type (
	LogLevelKey   int8
	LogLevelValue string
)

const (
	Debug LogLevelKey = iota + 1
	Info
	Warn
	Error
)

var LogLevelMap = map[LogLevelKey]LogLevelValue{
	Debug: "Debug",
	Info:  "Info",
	Warn:  "Warn",
	Error: "Error",
}

func GetLogLevel(key LogLevelKey) string {
	return string(LogLevelMap[key])
}

// TrackSource type in config
type TrackSource string

const (
	OFF TrackSource = "OFF"
	ON  TrackSource = "ON"
)
