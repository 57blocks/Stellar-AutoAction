package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/57blocks/auto-action/cli/internal/constant"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type (
	GlobalConfig struct {
		General   `toml:"general"`
		BoundWith `toml:"bound_with"`
	}
	GlobalConfigOpt func(sc *GlobalConfig)

	General struct {
		EnvPrefix string `toml:"env_prefix"`
		Log       string `toml:"log"`
	}

	BoundWith struct {
		Credential string `toml:"credential"`
		EndPoint   string `toml:"endpoint"`
	}
)

func Build(opts ...GlobalConfigOpt) *GlobalConfig {
	sc := new(GlobalConfig)

	for _, opt := range opts {
		opt(sc)
	}

	return sc
}

func WithLogLevel(logLevel string) GlobalConfigOpt {
	return func(sc *GlobalConfig) {
		sc.Log = logLevel
	}
}

func WithEnvPrefix(prefix string) GlobalConfigOpt {
	return func(sc *GlobalConfig) {
		sc.EnvPrefix = prefix
	}
}

func WithCredential(credential string) GlobalConfigOpt {
	return func(sc *GlobalConfig) {
		sc.Credential = credential
	}
}

func WithEndPoint(endpoint string) GlobalConfigOpt {
	return func(sc *GlobalConfig) {
		sc.EndPoint = endpoint
	}
}

// FindOrInit find or init the configuration file in the home directory
// together with the default configuration
func FindOrInit() (*GlobalConfig, string) {
	path := util.DefaultPath()

	if util.IsExists(path) {
		cfg, err := ReadConfig()
		cobra.CheckErr(err)

		return cfg, path
	}

	cfg := Build(
		WithCredential(util.DefaultCredPath()),
		WithEndPoint(constant.Host.String()),
		WithEnvPrefix(constant.EnvPrefix.ValStr()),
		WithLogLevel(constant.GetLogLevel(constant.Info)),
	)

	cobra.CheckErr(WriteConfig(cfg))

	return cfg, path
}

func ReadConfig() (*GlobalConfig, error) {
	data, err := os.ReadFile(util.DefaultPath())
	if err != nil {
		errMsg := fmt.Sprintf("read config error: %s\n", err.Error())
		slog.Error(errMsg)
		return nil, errors.Wrap(err, errMsg)
	}

	cfg := new(GlobalConfig)

	if _, err := toml.Decode(string(data), cfg); err != nil {
		errMsg := fmt.Sprintf("decode config error: %s\n", err.Error())
		slog.Error(errMsg)
		return nil, errors.Wrap(err, errMsg)
	}

	return cfg, nil
}

func WriteConfig(cfg *GlobalConfig) error {
	tomlBytes, err := toml.Marshal(cfg)
	if err != nil {
		errMsg := fmt.Sprintf("marshal config error: %s\n", err)
		slog.Error(errMsg)
		return errors.Wrap(err, errMsg)
	}

	if err := os.WriteFile(util.DefaultPath(), tomlBytes, 0666); err != nil {
		errMsg := fmt.Sprintf("write config error: %s\n", err.Error())
		slog.Error(errMsg)
		return errors.Wrap(err, errMsg)
	}

	return nil
}

func SyncConfigByFlags() error {
	cfg, err := ReadConfig()
	if err != nil {
		errMsg := fmt.Sprintf("read config error: %s\n", err.Error())
		slog.Error(errMsg)
		return errors.New(errMsg)
	}

	// Update fields if new values are provided
	if newCred := viper.GetString(constant.FlagCredential.ValStr()); newCred != "" {
		slog.Debug(fmt.Sprintf("newCredential: %v\n", newCred))
		cfg.Credential = newCred
	}
	if newEndPoint := viper.GetString(constant.FlagEndPoint.ValStr()); newEndPoint != "" {
		slog.Debug(fmt.Sprintf("newEndPoint: %v\n", newEndPoint))
		cfg.EndPoint = newEndPoint
	}
	if newEnvPrefix := viper.GetString(constant.FlagPrefix.ValStr()); newEnvPrefix != "" {
		slog.Debug(fmt.Sprintf("newEnvPrefix: %v\n", newEnvPrefix))
		cfg.EnvPrefix = newEnvPrefix
	}
	if newLogLevel := viper.GetString(constant.FlagLog.ValStr()); newLogLevel != "" {
		slog.Debug(fmt.Sprintf("newLogLevel: %v\n", newLogLevel))
		cfg.Log = newLogLevel
	}

	return WriteConfig(cfg)
}

func ResetConfigCredential() error {
	cfg, err := ReadConfig()
	if err != nil {
		errMsg := fmt.Sprintf("read config error: %s\n", err.Error())
		slog.Error(errMsg)
		return errors.New(errMsg)
	}

	cfg.Credential = ""

	return WriteConfig(cfg)
}
