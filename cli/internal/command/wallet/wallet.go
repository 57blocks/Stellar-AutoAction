package wallet

import (
	"github.com/57blocks/auto-action/cli/internal/command"

	"github.com/spf13/cobra"
)

var wallet = &cobra.Command{
	Use:   "wallet",
	Short: "Manage Stellar wallet addresses",
	Long: `
Description:
  The wallet command group provides tools for managing Stellar wallet addresses within the AutoAction system.
  It allows you to perform various operations related to your Stellar wallets.

Notes:
  - All operations are performed within the context of your user account.
  - Wallet addresses created and managed here are compatible with both Stellar mainnet and testnet.
  - Ensure you understand the implications of each action, especially when removing a wallet.

For detailed information on a specific subcommand, use:
  autoaction wallet <subcommand> --help
`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func init() {
	command.Root.AddCommand(wallet)
}
