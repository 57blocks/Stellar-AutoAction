package general

import (
	"github.com/57blocks/auto-action/cli/internal/command"

	"github.com/spf13/cobra"
)

// generalGroup represents the general command
var generalGroup = &cobra.Command{
	Use:   "general",
	Short: "CLI general settings",
}

func init() {
	command.Root.AddCommand(generalGroup)
}
