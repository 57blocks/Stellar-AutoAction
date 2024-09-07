package logx

import (
	"log/slog"

	"github.com/57blocks/auto-action/cli/internal/constant"
)

var Logger *slog.Logger

func SetupLogger(level string) {
	var logLevel slog.Level
	switch level {
	case constant.GetLogLevel(constant.Debug):
		//slog.SetLogLoggerLevel(slog.LevelDebug)
		logLevel = slog.LevelDebug
	case constant.GetLogLevel(constant.Warn):
		//slog.SetLogLoggerLevel(slog.LevelWarn)
		logLevel = slog.LevelWarn
	case constant.GetLogLevel(constant.Error):
		//slog.SetLogLoggerLevel(slog.LevelError)
		logLevel = slog.LevelError
	default:
		//slog.SetLogLoggerLevel(slog.LevelInfo)
		logLevel = slog.LevelInfo
	}

	opt := &slog.HandlerOptions{
		AddSource: false,
		Level:     logLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			return a
		},
	}

	Logger = slog.New(NewHandler(opt))
}
