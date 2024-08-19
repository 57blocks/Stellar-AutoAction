package config

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/spf13/viper"
)

var Global *Configuration

type (
	Configuration struct {
		Mode   string `mapstructure:"mode"`
		Bound  `mapstructure:"bound"`
		Log    `mapstructure:"log"`
		JWT    `mapstructure:"jwt"`
		Amazon `mapstructure:"amazon"`
	}

	Bound struct {
		Name     string `mapstructure:"name"`
		EndPoint string `mapstructure:"endpoint"`
	}

	Log struct {
		_        struct{}
		Level    string `mapstructure:"level"`
		Encoding string `mapstructure:"encoding"`
	}

	JWT struct {
		_          struct{}
		Protocol   string `mapstructure:"protocol"`
		PrivateKey string `mapstructure:"private_key"`
		PublicKey  string `mapstructure:"public_key"`
	}

	Amazon struct {
		_         struct{}
		Region    string `mapstructure:"region"`
		AccessKey string `mapstructure:"access_key"`
		SecretKey string `mapstructure:"secret_key"`
	}
)

func Setup() error {
	cfgLogger := slog.Default()
	slog.SetLogLoggerLevel(slog.LevelDebug)

	viper.NewWithOptions(
		viper.WithLogger(cfgLogger),
	)

	viper.AddConfigPath("./internal/config/")
	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.SetEnvPrefix("ST3LLAR")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	viper.AutomaticEnv()

	Global = new(Configuration)

	if err := viper.Unmarshal(&Global); err != nil {
		return err
	}

	cfgLogger.Debug(fmt.Sprintf("config path: %#v\n", viper.ConfigFileUsed()))
	cfgLogger.Debug(fmt.Sprintf("config: %#v\n", Global))

	return nil
}
