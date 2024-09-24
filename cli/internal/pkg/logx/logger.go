package logx

import (
	"log/slog"

	"github.com/57blocks/auto-action/cli/internal/constant"
)

var Logger *slog.Logger

func SetupLogger(level string, addSource bool) {
	var logLevel slog.Level
	switch level {
	case constant.GetLogLevel(constant.Debug):
		logLevel = slog.LevelDebug
	case constant.GetLogLevel(constant.Warn):
		logLevel = slog.LevelWarn
	case constant.GetLogLevel(constant.Error):
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opt := &slog.HandlerOptions{
		AddSource: addSource,
		Level:     logLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			return a
		},
	}

	Logger = slog.New(NewHandler(opt))
}
