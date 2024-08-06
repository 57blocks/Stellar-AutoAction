package config

import (
	"strings"

	"github.com/57blocks/auto-action/cli/internal/constant"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func SetupViper(cfg *GlobalConfig) {
	logger := SetupLogger(cfg) // TODO: abstraction layer for logger

	viper.NewWithOptions(
		viper.EnvKeyReplacer(strings.NewReplacer(".", "_")),
		viper.WithLogger(logger),
		//viper.KeyDelimiter("::"),
	)

	viper.AddConfigPath(util.Home())
	viper.SetConfigType(constant.ConfigurationType.ValStr())
	viper.SetConfigName(constant.ConfigName.ValStr())

	//viper.SetDefault("author", "v3nooom@outlook.com")
	//viper.SetDefault("license", "apache 2.0")

	viper.AutomaticEnv()

	cobra.CheckErr(viper.ReadInConfig())
}
