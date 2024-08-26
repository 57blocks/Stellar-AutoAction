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
		Amazon `mapstructure:"aws"`
		RDS    `mapstructure:"rds"`
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
		_               struct{}
		Region          string `mapstructure:"region"`
		AccessKeyID     string `mapstructure:"access_key_id"`
		SecretAccessKey string `mapstructure:"secret_access_key"`
	}

	RDS struct {
		_        struct{}
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Database string `mapstructure:"database"`
		SSLMode  string `mapstructure:"sslmode"`
	}

	Lambda struct {
		_    struct{}
		Role string `mapstructure:"role"`
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
