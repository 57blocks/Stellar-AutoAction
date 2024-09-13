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
	Short: "Create a wallet",
	Long: `
Description:	
  Create a new Stellar wallet address that can be used on
  both the mainnet and testnet of Stellar.
  However, before using the wallet,
  at least 1 XML must be transferred to this wallet address.

Note:
  - There is a limit to the number of wallet addresses that can be created. Currently, each user can only create 10 wallet addresses.
  - After creating the wallet address, you need to transfer at least 1 XML to this wallet address for it to function properly.
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
