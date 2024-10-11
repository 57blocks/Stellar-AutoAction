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
	Short: "Show the current authentication status",
	Long: `
Description:
  The status command displays the current authentication status for Stellar AutoAction.
  It provides a quick overview of your login state and associated details.

Output Information:
  - Currently authenticated account username
  - Associated organization name
  - Active Stellar Network

Behavior:
  - If you are logged in, the command will display the above information.
  - If you are not logged in, you will be prompted to authenticate first.

Example output:
  Account: johndoe, Organization: MyCompany, Network: Stellar Testnet

Related Commands:
  autoaction auth login  - Authenticate if not already logged in
  autoaction auth logout - End the current session

Note:
  This command is useful for verifying your current authentication state
  and ensuring you're operating in the correct context (account, organization, and network).
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
