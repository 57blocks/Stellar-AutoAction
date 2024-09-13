package oauth

import (
	"github.com/57blocks/auto-action/cli/internal/command"

	"github.com/spf13/cobra"
)

// oauthGroup represents the group command
var oauthGroup = &cobra.Command{
	Use:   "oauth",
	Short: "User oauth operations",
}

func init() {
	command.Root.AddCommand(oauthGroup)
}
