package log

import (
	"github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/constant"

	"go.uber.org/zap"
)

var Logger *ZapLogger

type (
	ZapFormatter struct{}
	ZapLogger    struct {
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

func (zl *ZapLogger) DEBUG(msg string, argMaps ...map[string]interface{}) {
	zl.ZapLog.Debug(msg, dispatcher(argMaps...)...)
}

func (zl *ZapLogger) INFO(msg string, argMaps ...map[string]interface{}) {
	zl.ZapLog.Info(msg, dispatcher(argMaps...)...)
}

func (zl *ZapLogger) WARN(msg string, argMaps ...map[string]interface{}) {
	zl.ZapLog.Warn(msg, dispatcher(argMaps...)...)
}

func (zl *ZapLogger) ERROR(msg string, argMaps ...map[string]interface{}) {
	zl.ZapLog.Error(msg, dispatcher(argMaps...)...)
}

func Setup() error {
	zapLogger, err := zap.Config{
		Level:             constant.GetLogLevel(config.Global.Level),
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: false,
		Encoding:          config.Global.Encoding,
		EncoderConfig:     zap.NewProductionEncoderConfig(),
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
		//InitialFields:     map[string]interface{}{"metadata": "metadata_sample"},
	}.Build()
	if err != nil {
		return err
	}

	Logger = &ZapLogger{zapLogger}
	Logger.DEBUG(
		"zap logger init",
		map[string]interface{}{"level": config.Global.Level, "encoding": config.Global.Encoding},
	)

	return nil
}
