package log

import (
	"github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/constant"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *ZapLoggerWrap

type (
	ZapFormatter  struct{}
	ZapLoggerWrap struct {
		ZapLog *zap.Logger
	}
)

var dispatcher = func(argMaps ...map[string]interface{}) []zap.Field {
	fs := make([]zap.Field, 0, len(argMaps))
	for _, argMap := range argMaps {
		for key, value := range argMap {
			switch v := value.(type) {
			case string:
				fs = append(fs, zap.String(key, v))
			case int:
				fs = append(fs, zap.Int(key, v))
			case bool:
				fs = append(fs, zap.Bool(key, v))
			case float64:
				fs = append(fs, zap.Float64(key, v))
			default:
				fs = append(fs, zap.Any(key, v))
			}
		}
	}

	return fs
}

func (zl *ZapLoggerWrap) DEBUG(msg string, argMaps ...map[string]interface{}) {
	zl.ZapLog.Debug(msg, dispatcher(argMaps...)...)
}

func (zl *ZapLoggerWrap) INFO(msg string, argMaps ...map[string]interface{}) {
	zl.ZapLog.Info(msg, dispatcher(argMaps...)...)
}

func (zl *ZapLoggerWrap) WARN(msg string, argMaps ...map[string]interface{}) {
	zl.ZapLog.Warn(msg, dispatcher(argMaps...)...)
}

func (zl *ZapLoggerWrap) ERROR(msg string, argMaps ...map[string]interface{}) {
	zl.ZapLog.Error(msg, dispatcher(argMaps...)...)
}

// Setup initializes the logger with configuration.
func Setup() error {
	zapLogger, err := zap.Config{
		Level:             constant.GetLogLevel(config.GlobalConfig.Level),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Encoding:          config.GlobalConfig.Encoding,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		//InitialFields:     map[string]interface{}{"metadata": "metadata_sample"},
	}.Build()
	if err != nil {
		return err
	}

	Logger = &ZapLoggerWrap{zapLogger}
	Logger.DEBUG(
		"zap logger init",
		map[string]interface{}{"level": config.GlobalConfig.Level, "encoding": config.GlobalConfig.Encoding},
	)

	return nil
}
