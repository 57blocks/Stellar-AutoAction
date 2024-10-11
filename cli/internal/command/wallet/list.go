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

var list = &cobra.Command{
	Use:   "list",
	Short: "Display all wallet addresses for the current user",
	Long: `
Description:
  The list command retrieves and displays all wallet addresses associated with the current user's account.
  This command provides a comprehensive overview of your Stellar wallets within the Stellar AutoAction system.

Output:
  The command returns a JSON-formatted list containing information about each wallet address.
  This may include details such as:
    - Wallet address (public key)

Note:
  - The list includes all wallets, regardless of their balance or activity status.
  - Ensure you are authenticated before running this command.

Related Commands:
  autoaction wallet create - Create a new wallet address
`,
	RunE: listFunc,
}

func init() {
	wallet.AddCommand(list)
}

func listFunc(_ *cobra.Command, _ []string) error {
	resp, err := supplierList()
	if err != nil {
		return err
	}

	var respData map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &respData); err != nil {
		logx.Logger.Error("Error unmarshalling JSON", "error", err.Error())
		return errorx.Internal(err.Error())
	}

	logx.Logger.Info("wallet list", "result", respData)

	return nil
}

func supplierList() (*resty.Response, error) {
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
		Get(URL)
	if err != nil {
		return nil, errorx.RestyError(err.Error())
	}
	if response.IsError() {
		return nil, errorx.WithRestyResp(response)
	}

	return response, nil
}
