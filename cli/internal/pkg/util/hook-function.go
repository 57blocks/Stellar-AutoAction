package util

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func PreBindFlags(cmd *cobra.Command, _ []string) {
	slog.Debug("PreBindFlags")

	if cmd.PersistentFlags().NFlag() > 0 {
		cmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
			viper.Set(flag.Name, flag.Value)
		})
	}

	if cmd.Flags().NFlag() > 0 {
		cmd.Flags().VisitAll(func(flag *pflag.Flag) {
			viper.Set(flag.Name, flag.Value)
		})
	}
}
