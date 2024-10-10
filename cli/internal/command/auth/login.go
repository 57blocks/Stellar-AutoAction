package auth

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/constant"
	"github.com/57blocks/auto-action/cli/internal/pkg/errorx"
	"github.com/57blocks/auto-action/cli/internal/pkg/logx"
	"github.com/57blocks/auto-action/cli/internal/pkg/restyx"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// login represents the login command
var login = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Stellar AutoAction",
	Long: `
Description:
  The login command authenticates you with the Stellar AutoAction system. It uses the credential 
  path specified in the configuration file by default. Alternatively, you can specify a custom 
  credential path using the --credential/-c flag.

Examples:
  autoaction auth login -a myaccount -o myorg
  autoaction auth login -c /path/to/custom/credential -a myaccount -o myorg

Behavior:
  - If it's your first time logging in or you're using a new credential path, the command will 
    create a new credential file at the specified location and update the configuration accordingly.
  - For subsequent logins, it will use the existing credential file.

Notes:
  - The account and organization flags are required for authentication.
  - You can manage multiple credentials using the 'configure' command.
  - Ensure you have the necessary permissions to create and modify credential files.
`,
	Args: cobra.NoArgs,
	RunE: loginFunc,
}

func init() {
	authGroup.AddCommand(login)

	flagCred := constant.FlagCredential.ValStr()
	login.Flags().StringP(
		flagCred,
		"c",
		config.Vp.GetString(flagCred),
		`Path to the credential file for authentication.
If omitted or first-time use, the default path will be used.
The command will create a new file if it doesn't exist.`)

	flagAcc := constant.FlagAccount.ValStr()
	login.Flags().StringP(
		flagAcc,
		"a",
		"",
		`Name of the account to authenticate.
This flag is required for login.`)

	flagOrg := constant.FlagOrganization.ValStr()
	login.Flags().StringP(flagOrg,
		"o",
		"",
		`Name of the organization associated with the account.
This flag is required for login.`)

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
	Password     string `json:"password"`
}

func loginFunc(cmd *cobra.Command, args []string) error {
	fmt.Println("Password: ")

	pwdBytes, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return errorx.Internal(fmt.Sprintf("reading passowrd error: %s", err.Error()))
	}

	if len(pwdBytes) == 0 {
		return errorx.BadRequest("empty password error")
	}

	key, err := util.LoadPublicKey(config.Vp.GetString("general.public_key"))
	if err != nil {
		return err
	}
	encodedPwd, err := util.EncryptPassword(string(pwdBytes), key)
	if err != nil {
		return err
	}

	success, err := supplierLogin(encodedPwd)
	if err != nil {
		return err
	}

	logx.Logger.Info("Login success! ")

	return syncLogin(success)
}

func supplierLogin(pwdHash string) (*resty.Response, error) {
	URL := util.ParseReqPath(fmt.Sprintf("%s/oauth/login", config.Vp.GetString("bound_with.endpoint")))

	response, err := restyx.Client.R().
		EnableTrace().
		SetBody(ReqLogin{
			Account:      config.Vp.GetString(constant.FlagAccount.ValStr()),
			Organization: config.Vp.GetString(constant.FlagOrganization.ValStr()),
			Password:     pwdHash,
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

func syncLogin(resp *resty.Response) error {
	cred := new(config.Credential)
	if err := json.Unmarshal(resp.Body(), cred); err != nil {
		return errorx.Internal(fmt.Sprintf("unmarshaling json response error: %s", err.Error()))
	}

	credPath := config.Vp.GetString(constant.FlagCredential.ValStr())
	if credPath == "" {
		credPath = util.DefaultCredPath()
	}

	if err := config.WriteCredential(credPath, cred); err != nil {
		return err
	}

	if err := config.ResetConfigCredential(credPath); err != nil {
		return err
	}

	return config.SyncConfigByFlags()
}
