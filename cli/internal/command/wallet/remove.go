package wallet

import (
	"fmt"

	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/pkg/errorx"
	"github.com/57blocks/auto-action/cli/internal/pkg/logx"
	"github.com/57blocks/auto-action/cli/internal/pkg/restyx"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"

	"github.com/spf13/cobra"
)

var remove = &cobra.Command{
	Use:   "remove [wallet-address]",
	Short: "Delete a wallet",
	Long: `
Description:	
  Delete a Stellar wallet address under a user.
  After the wallet address is deleted,
  it will no longer be able to sign third-party transactions,
  but the wallet address will still exist on the Stellar blockchain.

Note:
  - If the wallet address does not exist, the delete command will return an error.
  - You can only delete wallet addresses under your own user account.
`,
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
	if err := supplierRemove(walletAddress); err != nil {
		return err
	}

	logx.Logger.Info("remove wallet success")

	return nil
}

func supplierRemove(walletAddress string) error {
	token, err := config.Token()
	if err != nil {
		return err
	}

	URL := util.ParseReqPath(fmt.Sprintf("%s/wallet/%s", config.Vp.GetString("bound_with.endpoint"), walletAddress))

	response, err := restyx.Client.R().
		EnableTrace().
		SetHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": token,
		}).
		Delete(URL)
	if err != nil {
		return errorx.RestyError(err.Error())
	}
	if response.IsError() {
		return errorx.WithRestyResp(response)
	}

	return nil
}
