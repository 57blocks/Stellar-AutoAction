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

var create = &cobra.Command{
	Use:   "create",
	Short: "Generate a new Stellar wallet address",
	Long: `
Description:
  The create command generates a new Stellar wallet address that is compatible
  with both the Stellar mainnet and testnet. This wallet can be used for various
  Stellar network operations once activated.

Wallet Activation:
  To activate and use the newly created wallet, you must transfer at least 1 XLM
  (Stellar Lumens) to the generated address.

Important Notes:
  1. Wallet Limit: Each user is currently restricted to creating a maximum of 10 wallet addresses.
  2. Activation Requirement: The wallet will not be functional until it receives a minimum transfer of 1 XLM.
  3. Network Compatibility: The created address works on both Stellar mainnet and testnet.

Output:
  Upon successful creation, the command will display the new wallet address.

Example:
  autoaction wallet create

Next Steps:
  1. Securely store the generated wallet address.
  2. Transfer at least 1 XLM to the new address to activate it.
  3. Once activated, the wallet is ready for use on the Stellar network.
`,
	RunE: createFunc,
}

func init() {
	wallet.AddCommand(create)
}

func createFunc(_ *cobra.Command, _ []string) error {
	resp, err := supplierCreate()
	if err != nil {
		return err
	}

	wallet := make(map[string]interface{})
	if err := json.Unmarshal(resp.Body(), &wallet); err != nil {
		return errorx.Internal(fmt.Sprintf("unmarshaling json response error: %s", err.Error()))
	}

	logx.Logger.Info(fmt.Sprintf("create wallet success, address is %s", wallet["address"]))
	logx.Logger.Info("PS: Should deposit 1 XML to the new address to activate it.")

	return nil
}

func supplierCreate() (*resty.Response, error) {
	token, err := config.Token()
	if err != nil {
		logx.Logger.Error("PS: Should login first.")
		return nil, err
	}

	URL := util.ParseReqPath(fmt.Sprintf("%s/wallet", config.Vp.GetString("bound_with.endpoint")))

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
