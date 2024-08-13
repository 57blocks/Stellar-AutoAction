package constant

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	MappingLevel = map[LogLevel]zapcore.Level{
		LogDEBUG: zapcore.DebugLevel,
		LogINFO:  zapcore.InfoLevel,
		LogWARN:  zapcore.WarnLevel,
		LogERROR: zapcore.ErrorLevel,
		LogPANIC: zapcore.PanicLevel,
		LogFATAL: zapcore.FatalLevel,
	}
)

type (
	LogLevel string
)

const (
	LogDEBUG LogLevel = "DEBUG"
	LogINFO  LogLevel = "INFO"
	LogWARN  LogLevel = "WARN"
	LogERROR LogLevel = "ERROR"
	LogPANIC LogLevel = "PANIC"
	LogFATAL LogLevel = "FATAL"
)

func (l LogLevel) Val() string {
	return string(l)
}

func GetLogLevel(str string) zap.AtomicLevel {
	return zap.NewAtomicLevelAt(MappingLevel[LogLevel(str)])
}
