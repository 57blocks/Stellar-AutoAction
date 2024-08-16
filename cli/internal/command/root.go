package command

import (
	"fmt"
	"log/slog"

	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/pkg/errorx"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Root represents the base command when called without any subcommands
var (
	Root = &cobra.Command{
		Use:   "st3llar",
		Short: "The CLI toll for auto-action: st3llar",
		Long: `A CLI tool: st3llar, which helps users to run their 
method functions on Amazon Lambda quickly, together with the result
of the execution.`,
		Args:          cobra.OnlyValidArgs,
		ValidArgs:     []string{"configure", "login", "logout"}, // TODO: upgrade in-needed
		SilenceUsage:  true,
		SilenceErrors: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
			HiddenDefaultCmd:  true,
		},
		PersistentPreRun: util.PreBindFlags,
		RunE:             rootFunc,
	}
	version = "v0.0.1" // TODO: add release workflow to sync version with the git tag
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the Root.
func Execute() {
	errorx.CatchAndWrap(Root.Execute())
}

func init() {
	initConfig()

	Root.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
	Root.SetHelpCommand(&cobra.Command{
		Use:    "completion",
		Hidden: true,
	})
	Root.SetVersionTemplate(`Version: {{.Version}}`)
	Root.Version = version
}

func initConfig() {
	cfg, _ := config.FindOrInit()

	config.SetupViper(cfg)

	slog.Info(fmt.Sprintf("using config path: %s", viper.ConfigFileUsed()))
}

func rootFunc(cmd *cobra.Command, args []string) error {
	slog.Debug("---> rootFunc")
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		slog.Debug(fmt.Sprintf("flag.Name: %v, flag.Value: %v\n", flag.Name, flag.Value))
	})

	for _, v := range args {
		slog.Debug(fmt.Sprintf("args: %v\n", v))
	}

	return cmd.Usage()
}
