package wallet

import (
	"encoding/json"
	"fmt"

	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/pkg/errorx"
	"github.com/57blocks/auto-action/cli/internal/pkg/logx"
	"github.com/57blocks/auto-action/cli/internal/pkg/restyx"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

var verify = &cobra.Command{
	Use:   "verify [wallet-address]",
	Short: "Check the validity of a Stellar wallet address",
	Long: `
Description:
  The verify command allows you to check the validity of a Stellar wallet address associated with your user account.
  This command interacts with the Stellar blockchain to confirm the status of the specified wallet.

Arguments:
  [wallet-address]    The Stellar public key of the wallet you wish to verify

Verification Process:
  1. The command checks if the wallet address is associated with your user account.
  2. It then verifies the wallet's status on the Stellar blockchain.
  3. The wallet is considered valid if it has received any funds since its creation.

Output:
  The command will return the validity status of the wallet address.

Important Notes:
  - A wallet address is considered invalid if it has never received any funds after creation.
  - You can only verify wallet addresses associated with your own user account.
  - This command cannot be used to verify wallet addresses belonging to other users.

Example:
  autoaction wallet verify GXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX

Possible Results:
  - Valid: The wallet exists and has received funds.
  - Invalid: The wallet exists but has never received funds.
  - Error: The wallet is not associated with your account or doesn't exist.

Related Commands:
  autoaction wallet list - View all wallets in your account
  autoaction wallet create - Add a new wallet to your account
`,
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

	resp, err := supplierVerify(walletAddress)
	if err != nil {
		return err
	}

	wallet := make(map[string]interface{})
	if err := json.Unmarshal(resp.Body(), &wallet); err != nil {
		return errorx.Internal(fmt.Sprintf("unmarshaling json response error: %s", err.Error()))
	}
	message := fmt.Sprintf("The wallet address %s is VALID", walletAddress)
	if !wallet["is_valid"].(bool) {
		message = fmt.Sprintf(`The wallet address %s is INVALID. 
It might be because the wallet needs at least 1 XML to activate, 
or the Stellar has cleared the address(Testnet only).`, walletAddress)
	}
	logx.Logger.Info(message)

	return nil
}

func supplierVerify(walletAddress string) (*resty.Response, error) {
	token, err := config.Token()
	if err != nil {
		return nil, err
	}

	URL := util.ParseReqPath(fmt.Sprintf("%s/wallet/%s", config.Vp.GetString("bound_with.endpoint"), walletAddress))

	response, err := restyx.Client.R().
		EnableTrace().
		SetHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": token,
		}).
		Post(URL)
	if err != nil {
		return nil, errorx.RestyError(err.Error())
	}
	if response.IsError() {
		return nil, errorx.WithRestyResp(response)
	}

	return response, nil
}
