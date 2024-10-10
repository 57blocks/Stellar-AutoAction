package auth

import (
	"github.com/57blocks/auto-action/cli/internal/command"

	"github.com/spf13/cobra"
)

// authGroup represents the group command
var authGroup = &cobra.Command{
	Use:   "auth",
	Short: "Manage user authentication",
	Long: `
Description:
  The auth command group provides various subcommands for managing user authentication 
  in the AutoAction system.
`,
}

func init() {
	command.Root.AddCommand(authGroup)
}
