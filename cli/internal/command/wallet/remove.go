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
	Short: "Delete a Stellar wallet from your account",
	Long: `
Description:
  The remove command allows you to delete a specific Stellar wallet address associated with your user account.
  This action removes the wallet from the AutoAction system, but does not affect the wallet's existence on the
  Stellar blockchain.

Arguments:
  [wallet-address]    The Stellar public key of the wallet you wish to remove

Effects of Removal:
  1. The wallet will be disassociated from your AutoAction account.
  2. The wallet will no longer be able to sign third-party transactions through AutoAction.
  3. The wallet and its funds will continue to exist on the Stellar blockchain.

Important Notes:
  - You can only remove wallet addresses associated with your own AutoAction account.
  - If the specified wallet address does not exist in your account, the command will return an error.
  - This action cannot be undone. Make sure you want to remove the wallet before proceeding.

Example:
  autoaction wallet remove GXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX

Caution:
  Removing a wallet from AutoAction does not delete it from the Stellar network. Ensure you have
  access to the wallet's secret key if you intend to use it outside of AutoAction in the future.

Related Commands:
  autoaction wallet list - View all wallets in your account
  autoaction wallet create - Add a new wallet to your account
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
