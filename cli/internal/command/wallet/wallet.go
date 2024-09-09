package wallet

import (
	"fmt"

	"github.com/57blocks/auto-action/cli/internal/command"

	"github.com/spf13/cobra"
)

var wallet = &cobra.Command{
	Use:   "wallet",
	Short: "Wallet Address Management",
	Long: `
Description:
  The wallet address management command is used to 
  create, delete, verify, and list wallet addresses.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("wallet called")
	},
}

func init() {
	command.Root.AddCommand(wallet)
}
