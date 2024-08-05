package oauth

import (
	"github.com/57blocks/auto-action/cli/internal/command"

	"github.com/spf13/cobra"
)

// logout represents the logout command
var logout = &cobra.Command{
	Use:   "logout",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely
contains examples and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
