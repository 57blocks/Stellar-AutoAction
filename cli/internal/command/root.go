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
		Short: "CLI tool for managing and executing actions on AutoAction",
		Long: `
AutoAction: A Comprehensive CLI Tool for Action Management and Execution

Description:
  AutoAction is a powerful command-line interface (CLI) tool designed to streamline 
  the process of managing and executing actions within the AutoAction ecosystem. 
  It offers a range of functionalities to enhance your workflow and productivity.

Key Features:
  1. Action Execution: Quickly run your handlers on AutoAction.
  2. Log Tracking: Monitor and track execution logs for your actions.
  3. Wallet Management: Create and manage wallets using CubeSigner integration.
  4. Transaction Signing: Securely sign transactions for blockchain operations.

For more information about a specific command, use:
  autoaction [command] --help

Note:
  Ensure you have the necessary permissions and are authenticated 
  before performing sensitive operations.
`,
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
