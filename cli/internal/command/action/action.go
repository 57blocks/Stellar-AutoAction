package action

import (
	"github.com/57blocks/auto-action/cli/internal/command"

	"github.com/spf13/cobra"
)

// actionGroup represents the action command
var actionGroup = &cobra.Command{
	Use:   "action",
	Short: "Action management",
}

func init() {
	command.Root.AddCommand(actionGroup)
}
