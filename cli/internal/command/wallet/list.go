package wallet

import (
	"github.com/57blocks/auto-action/cli/internal/pkg/logx"

	"github.com/spf13/cobra"
)

var list = &cobra.Command{
	Use:   "list",
	Short: "List wallets",
	Long: `
Description:	
  List all saved wallet addresses.`,
	RunE: listFunc,
}

func init() {
	wallet.AddCommand(list)
}

func listFunc(_ *cobra.Command, _ []string) error {
	logx.Logger.Info("list called")
	return nil
}
