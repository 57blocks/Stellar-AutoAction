package wallet

import (
	"encoding/json"
	"fmt"

	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/constant"
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

	flagEnv := constant.FlagEnv.ValStr()
	verify.Flags().StringP(
		flagEnv,
		"e",
		config.Vp.GetString(flagEnv),
		`
The environment of the wallet address. 
The value should be one of the following:
- testnet (default)
- mainnet
`)
}

func verifyFunc(_ *cobra.Command, args []string) error {
	walletAddress := args[0]
	env := config.Vp.GetString(constant.FlagEnv.ValStr())
	if env == "" {
		env = "testnet"
	}

	logx.Logger.Info(fmt.Sprintf("Verifying wallet with address: %s\n", walletAddress))

	resp, err := supplierVerify(walletAddress, env)
	if err != nil {
		return errorx.Internal(fmt.Sprintf("verify wallet failed: %s", err.Error()))
	}

	wallet := new(config.Wallet)
	if err := json.Unmarshal(resp.Body(), wallet); err != nil {
		return errorx.Internal(fmt.Sprintf("unmarshaling json response error: %s", err.Error()))
	}
	isValid := "invalid"
	if wallet.IsValid {
		isValid = "valid"
	}
	logx.Logger.Info(fmt.Sprintf("The wallet adderss %s in %s is %s", walletAddress, env, isValid))

	return nil
}

func supplierVerify(walletAddress string, env string) (*resty.Response, error) {
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
		SetBody(map[string]string{
			"env": env,
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
