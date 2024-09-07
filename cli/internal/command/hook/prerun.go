package hook

import (
	"log/slog"

	"github.com/57blocks/auto-action/cli/internal/config"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func PreRunFunc(cmd *cobra.Command, _ []string) {
	slog.Debug("Root commandPreRunFunc")

	if cmd.PersistentFlags().NFlag() > 0 {
		cmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
			config.Vp.Set(flag.Name, flag.Value)
		})
	}

	if cmd.Flags().NFlag() > 0 {
		cmd.Flags().VisitAll(func(flag *pflag.Flag) {
			config.Vp.Set(flag.Name, flag.Value)
		})
	}
}
