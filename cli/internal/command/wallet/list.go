package wallet

import (
	"fmt"

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
	fmt.Println("list called")
	return nil
}
