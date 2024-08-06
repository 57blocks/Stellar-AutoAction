package util

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func PreRunBindFlags(cmd *cobra.Command, _ []string) {
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
