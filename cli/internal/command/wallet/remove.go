package wallet

import (
	"fmt"

	"github.com/57blocks/auto-action/cli/internal/pkg/errorx"
	"github.com/57blocks/auto-action/cli/internal/pkg/logx"

	"github.com/spf13/cobra"
)

var remove = &cobra.Command{
	Use:   "remove [wallet-address]",
	Short: "Delete a wallet",
	Long: `
Description:	
  Delete an existing wallet address.`,
	Args: cobra.ExactArgs(1),
	RunE: removeFunc,
}

func init() {
	wallet.AddCommand(remove)

	remove.SetUsageTemplate(`
Usage:
  wallet remove [wallet-address]

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}`)
}

func removeFunc(_ *cobra.Command, args []string) error {
	walletAddress := args[0]
	logx.Logger.Info(fmt.Sprintf("Removing wallet with address: %s\n", walletAddress))
	if true {
		logx.Logger.Info("remove wallet success")
		return nil
	}
	logx.Logger.Info("remove wallet failed")
	return errorx.BadRequest("remove wallet failed")
}
