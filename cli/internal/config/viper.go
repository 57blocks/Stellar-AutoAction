package config

import (
	"strings"

	"github.com/57blocks/auto-action/cli/internal/constant"
	"github.com/57blocks/auto-action/cli/internal/pkg/logx"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Vp *viper.Viper
)

func SetupViper(cfg *GlobalConfig) {
	trackSource := false
	if strings.EqualFold(cfg.Source, "on") {
		trackSource = true
	}

	logx.SetupLogger(cfg.Log, trackSource)

	Vp = viper.NewWithOptions(
		viper.EnvKeyReplacer(strings.NewReplacer(".", "_")),
		viper.WithLogger(logx.Logger),
	)

	Vp.AddConfigPath(util.Home())
	Vp.SetConfigType(constant.ConfigurationType.ValStr())
	Vp.SetConfigName(constant.ConfigName.ValStr())

	Vp.AutomaticEnv()

	cobra.CheckErr(Vp.ReadInConfig())
}
