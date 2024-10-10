package auth

import (
	"fmt"

	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/pkg/errorx"

	"github.com/spf13/cobra"
)

// status represents the status command
var status = &cobra.Command{
	Use:   "status",
	Short: "Display the login status of auto-action.",
	Long: `
Description:
  Display the current login status of auto-action.

This command will show the following information:
  - Currently logged-in account
  - Associated organization
  - Stellar Network in use

If you are not logged in, this command will prompt you to log in first.
`,
	Args: cobra.NoArgs,
	RunE: statusFunc,
}

func init() {
	authGroup.AddCommand(status)
}

func statusFunc(cmd *cobra.Command, args []string) error {
	cfg, err := config.ReadConfig()
	if err != nil {
		return errorx.Internal(fmt.Sprintf("show status failed, you should login first. error: %s", err.Error()))
	}

	cred, err := config.ReadCredential(cfg.BoundWith.Credential)
	if err != nil {
		return errorx.Internal(fmt.Sprintf("show status failed, you should login first. error: %s", err.Error()))
	}

	fmt.Printf("Account: %s, Organization: %s, Network: %s\n", cred.Account, cred.Organization, cred.Network)
	return nil
}
