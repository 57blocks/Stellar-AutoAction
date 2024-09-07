package oauth

import (
	"encoding/json"
	"fmt"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"
	"os"

	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/constant"
	"github.com/57blocks/auto-action/cli/internal/pkg/errorx"
	"github.com/57blocks/auto-action/cli/internal/pkg/restyx"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// login represents the login command
var login = &cobra.Command{
	Use:   "login",
	Short: "Login to the Stellar auto-action.",
	Long: `
Description:
  Login the Stellar auto-action based on credential path in the config.
  Or, by the --credential/-c flag, to specify the credential path.
  And will create a new credential under the path you just
  claimed and set it to config, when it's the first time.

Note:
  - You could specify other credentials by **configure** command.
`,
	Args: cobra.NoArgs,
	RunE: loginFunc,
}

func init() {
	oauthGroup.AddCommand(login)

	flagCred := constant.FlagCredential.ValStr()
	login.Flags().StringP(
		flagCred,
		"c",
		config.Vp.GetString(flagCred),
		`
The credential file for the command about to be bound.
If it's the first time, or ignored, the default path will be used.'`)

	flagEnv := constant.FlagEnvironment.ValStr()
	login.Flags().StringP(
		flagEnv,
		"e",
		config.Vp.GetString(flagEnv),
		`
The execution environment about to be bound.
If ignored, the default environment: Horizon will be used.'`)

	flagAcc := constant.FlagAccount.ValStr()
	login.Flags().StringP(
		flagAcc,
		"a",
		"",
		"name of the account")

	flagOrg := constant.FlagOrganization.ValStr()
	login.Flags().StringP(flagOrg,
		"o",
		"",
		"name of the organization")

	if err := login.MarkFlagRequired(flagAcc); err != nil {
		return
	}
	if err := login.MarkFlagRequired(flagOrg); err != nil {
		return
	}
}

type ReqLogin struct {
	Account      string `json:"account"`
	Organization string `json:"organization"`
	Password     []byte `json:"password"`
	Environment  string `json:"environment"`
}

func loginFunc(cmd *cobra.Command, args []string) error {
	fmt.Println("Password: ")

	pwdBytes, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return errorx.Internal(fmt.Sprintf("reading passowrd error: %s", err.Error()))
	}

	if len(pwdBytes) == 0 {
		return errorx.BadRequest("empty cryptPwd error")
	}

	success, err := supplierLogin(pwdBytes)
	if err != nil {
		return err
	}

	return syncLogin(success)
}

func supplierLogin(cryptPwdBytes []byte) (*resty.Response, error) {
	URL := util.ParseReqPath(fmt.Sprintf("%s/oauth/login", config.Vp.GetString("bound_with.endpoint")))

	response, err := restyx.Client.R().
		EnableTrace().
		SetBody(ReqLogin{
			Account:      config.Vp.GetString(constant.FlagAccount.ValStr()),
			Organization: config.Vp.GetString(constant.FlagOrganization.ValStr()),
			Password:     cryptPwdBytes,
			Environment:  config.Vp.GetString(constant.FlagEnvironment.ValStr()),
		}).
		Post(URL)
	if err != nil {
		return nil, errorx.WithRestyResp(response)
	}
	if response.IsError() {
		return nil, errorx.WithRestyResp(response)
	}

	return response, nil
}

func syncLogin(resp *resty.Response) error {
	cred := new(config.Credential)
	if err := json.Unmarshal(resp.Body(), cred); err != nil {
		return errorx.Internal(fmt.Sprintf("unmarshaling json response error: %s", err.Error()))
	}

	if err := config.WriteCredential(config.Vp.GetString(constant.FlagCredential.ValStr()), cred); err != nil {
		return err
	}

	return config.SyncConfigByFlags()
}
