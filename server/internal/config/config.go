package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

var (
	GlobalConfig *Configuration
	Vp           *viper.Viper
)

type (
	Configuration struct {
		Mode   string `mapstructure:"mode"`
		Bound  `mapstructure:"bound"`
		Log    `mapstructure:"log"`
		PEM    `mapstructure:"pem"`
		JWT    `mapstructure:"jwt"`
		Amazon `mapstructure:"aws"`
		RDS    `mapstructure:"rds"`
		CS     `mapstructure:"cs"`
		Wallet `mapstructure:"wallet"`
		Lambda `mapstructure:"lambda"`
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

	PEM struct {
		_          struct{}
		PrivateKey string `mapstructure:"private_key"`
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

	CS struct {
		_            struct{}
		Endpoint     string `mapstructure:"endpoint"`
		Organization string `mapstructure:"organization"`
	}

	Wallet struct {
		_   struct{}
		Max int `mapstructure:"max"`
	}

	Lambda struct {
		_   struct{}
		Max int `mapstructure:"max"`
	}
)

func Setup() error {
	cfgLogger := slog.Default()
	slog.SetLogLoggerLevel(slog.LevelDebug)

	Vp = viper.NewWithOptions(
		viper.WithLogger(cfgLogger),
	)

	Vp.AddConfigPath("./internal/config/")
	Vp.SetConfigType("toml")
	Vp.SetConfigName("config")
	Vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := Vp.ReadInConfig(); err != nil {
		return err
	}

	Vp.AutomaticEnv()

	GlobalConfig = new(Configuration)

	if err := Vp.Unmarshal(&GlobalConfig); err != nil {
		return err
	}

	cfgLogger.Debug(fmt.Sprintf("config path: %#v", Vp.ConfigFileUsed()))
	cfgLogger.Debug(fmt.Sprintf("config: %#v", GlobalConfig.DebugStr()))

	return nil
}

func (c *Configuration) DebugStr() string {
	debugMap := make(map[string]interface{})
	debugMapRecursive(reflect.ValueOf(*c), "", debugMap)
	jsonBytes, _ := json.Marshal(debugMap)
	return string(jsonBytes)
}

func debugMapRecursive(v reflect.Value, prefix string, debugMap map[string]interface{}) {
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		fieldName := prefix + field.Name

		switch field.Type.Kind() {
		case reflect.String:
			strValue := value.String()
			if len(strValue) > 20 {
				strValue = strValue[:20] + "..."
			}
			debugMap[fieldName] = strValue
		case reflect.Struct:
			subMap := make(map[string]interface{})
			debugMap[fieldName] = subMap
			debugMapRecursive(value, fieldName+".", subMap)
		default:
			debugMap[fieldName] = value.Interface()
		}
	}
}
