package action

import (
	"github.com/57blocks/auto-action/cli/internal/command"

	"github.com/spf13/cobra"
)

// actionGroup represents the action command
var actionGroup = &cobra.Command{
	Use:   "action",
	Short: "Manage actions in Stellar AutoAction",
	Long: `
Description:
  The action command provides tools for managing actions within Stellar AutoAction.

This command group allows you to:
  - List existing actions
  - Create new actions
  - Update action configurations
  - Delete actions
  - Execute actions
  - View action execution history
`,
}

func init() {
	command.Root.AddCommand(actionGroup)
}
