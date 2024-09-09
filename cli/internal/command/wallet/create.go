package wallet

import (
	"github.com/57blocks/auto-action/cli/internal/pkg/logx"

	"github.com/spf13/cobra"
)

var create = &cobra.Command{
	Use:   "create",
	Short: "Create a wallet",
	Long: `
Description:	
  Create a new wallet address.`,
	RunE: createFunc,
}

func init() {
	wallet.AddCommand(create)
}

func createFunc(_ *cobra.Command, _ []string) error {
	logx.Logger.Info("create wallet success, address: 0x1234567890123456789012345678901234567890")
	logx.Logger.Info("PS: Should deposit 1 XML to the new address to activate it.")
	return nil
}
