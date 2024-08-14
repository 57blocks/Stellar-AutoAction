package config

import (
	"log/slog"

	"github.com/57blocks/auto-action/cli/internal/constant"
)

func SetupLogger(cfg *GlobalConfig) *slog.Logger {
	logger := slog.Default()

	switch cfg.Log {
	case constant.GetLogLevel(constant.Debug):
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case constant.GetLogLevel(constant.Warn):
		slog.SetLogLoggerLevel(slog.LevelWarn)
	case constant.GetLogLevel(constant.Info):
		slog.SetLogLoggerLevel(slog.LevelError)
	default:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}

	return logger
}
