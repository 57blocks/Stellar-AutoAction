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
		Logger `mapstructure:"logger"`
		JWT    `mapstructure:"jwt"`
		Amazon `mapstructure:"amazon"`
	}

	Logger struct {
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
		//viper.EnvKeyReplacer(strings.NewReplacer(".", "_")),
		viper.WithLogger(cfgLogger),
	)

	//viper.SetEnvPrefix("ST3LLAR")
	viper.AddConfigPath("./internal/config/")
	viper.SetConfigType("toml")
	viper.SetConfigName("config")

	viper.SetEnvPrefix("")
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

	//fmt.Printf("jwt.protocol: %s\n", viper.GetString("jwt.protocol"))
	//fmt.Printf("jwt.private : %s\n", viper.GetString("jwt.private"))
	//fmt.Printf("jwt.public  : %s\n", viper.GetString("jwt.public"))
	//
	//fmt.Printf("amazon.region: %s\n", viper.GetString("amazon.region"))
	//fmt.Printf("amazon.access_key: %s\n", viper.GetString("amazon.access_key"))
	//fmt.Printf("amazon.secret_key: %s\n", viper.GetString("amazon.secret_key"))

	return nil
}
