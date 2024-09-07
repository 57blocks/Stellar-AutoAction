package wallet

import (
	"errors"
	"fmt"

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
	fmt.Printf("Removing wallet with address: %s\n", walletAddress)
	if true {
		fmt.Println("remove wallet success")
		return nil
	}
	fmt.Println("remove wallet failed")
	return errors.New("remove wallet failed")
}
