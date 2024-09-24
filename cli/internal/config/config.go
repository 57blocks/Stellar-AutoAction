package config

import (
	"fmt"
	"os"

	"github.com/57blocks/auto-action/cli/internal/constant"
	"github.com/57blocks/auto-action/cli/internal/pkg/errorx"
	"github.com/57blocks/auto-action/cli/internal/pkg/logx"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

type (
	GlobalConfig struct {
		General   `toml:"general"`
		BoundWith `toml:"bound_with"`
	}
	GlobalConfigOpt func(sc *GlobalConfig)

	General struct {
		Log       string `toml:"logx"`
		Source    string `toml:"source"`
		PublicKey string `toml:"public_key"`
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

func WithTrackSource(source string) GlobalConfigOpt {
	return func(sc *GlobalConfig) {
		sc.Source = source
	}
}

func WithPublicKey(pubKey string) GlobalConfigOpt {
	return func(sc *GlobalConfig) {
		sc.PublicKey = pubKey
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
		WithEndPoint(""),
		WithLogLevel(constant.GetLogLevel(constant.Info)),
		WithTrackSource(string(constant.OFF)),
		WithPublicKey(""),
	)

	cobra.CheckErr(WriteConfig(cfg))

	return cfg, path
}

func ReadConfig() (*GlobalConfig, error) {
	data, err := os.ReadFile(util.DefaultPath())
	if err != nil {
		return nil, errorx.Internal(fmt.Sprintf("read config error: %s", err.Error()))
	}

	cfg := new(GlobalConfig)

	if _, err := toml.Decode(string(data), cfg); err != nil {
		return nil, errorx.Internal(fmt.Sprintf("decode config error: %s", err.Error()))
	}

	return cfg, nil
}

func WriteConfig(cfg *GlobalConfig) error {
	tomlBytes, err := toml.Marshal(cfg)
	if err != nil {
		return errorx.Internal(fmt.Sprintf("marshal config error: %s", err))
	}

	if err := os.WriteFile(util.DefaultPath(), tomlBytes, 0666); err != nil {
		return errorx.Internal(fmt.Sprintf("write config error: %s", err.Error()))
	}

	return nil
}

func SyncConfigByFlags() error {
	cfg, err := ReadConfig()
	if err != nil {
		return errorx.Internal(fmt.Sprintf("read config error: %s", err.Error()))
	}

	// Update fields if new values are provided
	if newCred := Vp.GetString(constant.FlagCredential.ValStr()); newCred != "" {
		logx.Logger.Debug("sync credential", "updated to", newCred)
		cfg.Credential = newCred
	}
	if newEndPoint := Vp.GetString(constant.FlagEndPoint.ValStr()); newEndPoint != "" {
		logx.Logger.Debug("sync endpoint", "updated to", newEndPoint)
		cfg.EndPoint = newEndPoint
	}
	if newLogLevel := Vp.GetString(constant.FlagLog.ValStr()); newLogLevel != "" {
		logx.Logger.Debug("sync logx level", "updated to", newLogLevel)
		cfg.Log = newLogLevel
	}
	if newSource := Vp.GetString(constant.FlagSource.ValStr()); newSource == string(constant.ON) ||
		newSource == string(constant.OFF) {
		logx.Logger.Debug("sync tracking source or not", "updated to", newSource)
		cfg.Source = newSource
	}

	return WriteConfig(cfg)
}

func ResetConfigCredential() error {
	cfg, err := ReadConfig()
	if err != nil {
		return errorx.Internal(fmt.Sprintf("read config error: %s", err.Error()))
	}

	cfg.Credential = ""

	return WriteConfig(cfg)
}
