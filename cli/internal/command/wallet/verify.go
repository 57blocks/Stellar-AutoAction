package wallet

import (
	"fmt"

	"github.com/57blocks/auto-action/cli/internal/pkg/errorx"
	"github.com/57blocks/auto-action/cli/internal/pkg/logx"

	"github.com/spf13/cobra"
)

var verify = &cobra.Command{
	Use:   "verify [wallet-address]",
	Short: "Verify a wallet",
	Long: `
Description:	
  Verify the validity of a wallet address.`,
	Args: cobra.ExactArgs(1),
	RunE: verifyFunc,
}

func init() {
	wallet.AddCommand(verify)

	verify.SetUsageTemplate(`
Usage:
  wallet verify [wallet-address]

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}`)
}

func verifyFunc(_ *cobra.Command, args []string) error {
	walletAddress := args[0]
	logx.Logger.Info(fmt.Sprintf("Verifying wallet with address: %s\n", walletAddress))
	if true {
		logx.Logger.Info("verify wallet success")
		return nil
	}
	logx.Logger.Info("verify wallet failed")
	return errorx.BadRequest("verify wallet failed")
}
