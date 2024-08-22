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
		cfg, err := ReadConfig(path)
		cobra.CheckErr(err)

		return cfg, path
	}

	cfg := Build(
		WithCredential(util.DefaultCredPath()),
		WithEndPoint(constant.Host.String()),
		WithEnvPrefix(constant.EnvPrefix.ValStr()),
		WithLogLevel(constant.GetLogLevel(constant.Info)),
	)

	cobra.CheckErr(WriteConfig(cfg, path))

	return cfg, path
}

func ReadConfig(path string) (*GlobalConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := new(GlobalConfig)

	if _, err := toml.Decode(string(data), cfg); err != nil {
		_, e := fmt.Fprintf(
			os.Stderr,
			"reading config error: %s\n",
			err.Error(),
		)
		return nil, e
	}

	return cfg, nil
}

func WriteConfig(cfg *GlobalConfig, path string) error {
	tomlBytes, err := toml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshalling config error: %w", err)
	}

	if err := os.WriteFile(path, tomlBytes, 0666); err != nil {
		_, e := fmt.Fprintf(
			os.Stderr,
			"writing config error: %s\n",
			err.Error(),
		)
		return e
	}

	return nil
}

func SyncConfigByFlags(path string) error {
	cfg, err := ReadConfig(util.DefaultPath())
	if err != nil {
		return errors.New(fmt.Sprintf("reading configuration error: %s\n", err.Error()))
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

	return WriteConfig(cfg, path)
}
