package command

import (
	"github.com/57blocks/auto-action/cli/internal/command/hook"
	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/pkg/errorx"
	"github.com/57blocks/auto-action/cli/internal/pkg/logx"

	"github.com/spf13/cobra"
)

// Root represents the base command when called without any subcommands
var (
	Root = &cobra.Command{
		Use:   "autoaction",
		Short: "The CLI toll for auto-action: autoaction",
		Long: `
A CLI tool: autoaction, which helps users to run their handler on AutoAction 
quickly, together with tracking the execution logs.

Also, it provides a way to create wallet through the CubeSigner and sign
the transactions.
`,
		Args:          cobra.OnlyValidArgs,
		ValidArgs:     []string{"configure", "login", "logout"}, // TODO: upgrade in-needed
		SilenceUsage:  true,
		SilenceErrors: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
			HiddenDefaultCmd:  true,
		},
		PersistentPreRun: hook.PreRunFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
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

	logx.Logger.Info("init success", "using config path", config.Vp.ConfigFileUsed())
}
