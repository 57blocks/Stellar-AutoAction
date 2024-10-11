package general

import (
	"github.com/57blocks/auto-action/cli/internal/command"

	"github.com/spf13/cobra"
)

// generalGroup represents the general command
var generalGroup = &cobra.Command{
	Use:   "general",
	Short: "Manage general settings for Stellar AutoAction CLI",
	Long: `
Description:
  The general command group provides tools to configure and manage various
  settings for the Stellar AutoAction CLI. These settings affect the overall behavior
  and functionality of the CLI.

Notes:
  - Changes made using these commands will affect all future operations of the CLI.
  - It's recommended to review your settings periodically to ensure optimal performance.
`,
}

func init() {
	command.Root.AddCommand(generalGroup)
}
