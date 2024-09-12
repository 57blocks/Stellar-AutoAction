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
	Short: "Verify a wallet",
	// TODO: add long description
	Long: `
Description:	
  Verify the validity of a wallet address.
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
	isValid := "invalid"
	if wallet["is_valid"] == true {
		isValid = "valid"
	}
	logx.Logger.Info(fmt.Sprintf("The wallet adderss %s is %s", walletAddress, isValid))

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
