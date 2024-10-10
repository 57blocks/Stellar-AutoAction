package auth

import (
	"github.com/57blocks/auto-action/cli/internal/command"

	"github.com/spf13/cobra"
)

// authGroup represents the group command
var authGroup = &cobra.Command{
	Use:   "auth",
	Short: "User authentication operations",
}

func init() {
	command.Root.AddCommand(authGroup)
}
