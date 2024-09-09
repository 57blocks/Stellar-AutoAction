package wallet

import (
	"errors"
	"fmt"

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

	remove.SetUsageTemplate(`
Usage:
  wallet verify [wallet-address]

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}`)
}

func verifyFunc(_ *cobra.Command, args []string) error {
	walletAddress := args[0]
	fmt.Printf("Verifying wallet with address: %s\n", walletAddress)
	if true {
		fmt.Println("verify wallet success")
		return nil
	}
	fmt.Println("verify wallet failed")
	return errors.New("verify wallet failed")
}
