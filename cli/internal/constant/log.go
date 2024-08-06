package constant

type (
	LogLevelKey   int8
	LogLevelValue string
)

var LogLevelMap = map[LogLevelKey]LogLevelValue{
	Debug: "Debug",
	Info:  "Info",
	Warn:  "Warn",
	Error: "Error",
}

const (
	Debug LogLevelKey = iota + 1
	Info
	Warn
	Error
)

func GetLogLevel(key LogLevelKey) string {
	return string(LogLevelMap[key])
}
