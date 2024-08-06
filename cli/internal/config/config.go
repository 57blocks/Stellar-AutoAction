package config

import (
	"fmt"
	"os"

	"github.com/57blocks/auto-action/cli/internal/constant"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

type (
	GlobalConfig struct {
		Credential   string `toml:"credential"`
		EnvVarPrefix string `toml:"env-var-prefix"`
		LogLevel     string `toml:"log-level"`
	}
	GlobalCfgOpt func(sc *GlobalConfig)
)

func Build(opts ...GlobalCfgOpt) *GlobalConfig {
	sc := new(GlobalConfig)

	for _, opt := range opts {
		opt(sc)
	}

	return sc
}

func WithLogLevel(logLevel string) GlobalCfgOpt {
	return func(sc *GlobalConfig) {
		sc.LogLevel = logLevel
	}
}

func WithEnvVarPrefix(prefix string) GlobalCfgOpt {
	return func(sc *GlobalConfig) {
		sc.EnvVarPrefix = prefix
	}
}

func WithCredential(credential string) GlobalCfgOpt {
	return func(sc *GlobalConfig) {
		sc.Credential = credential
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
		WithEnvVarPrefix(constant.EnvVarPrefix.ValStr()),
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
