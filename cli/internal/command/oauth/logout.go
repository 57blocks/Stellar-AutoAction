package oauth

import (
	"github.com/57blocks/auto-action/cli/internal/command"

	"github.com/spf13/cobra"
)

// logout represents the logout command
var logout = &cobra.Command{
	Use:   "logout",
	Short: "Logout the current session",
	Long: `Logout the current session under the path of credential in
the config.

For other credentials, they are still alive. It's recommended that
switching session by using **configure** command to set other 
credentials to the config`,
	Args: cobra.NoArgs,
	RunE: logoutFunc,
}

func init() {
	command.Root.AddCommand(logout)
}

func logoutFunc(_ *cobra.Command, _ []string) error {
	// TODO: add implementation

	return nil
}
