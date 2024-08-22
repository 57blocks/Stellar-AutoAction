package oauth

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/57blocks/auto-action/cli/internal/command"
	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/constant"
	"github.com/57blocks/auto-action/cli/internal/pkg/restyx"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

// login represents the login command
var login = &cobra.Command{
	Use:   "login",
	Short: "Login to the Stellar auto-action.",
	Long: `Login the Stellar auto-action based on credential path in
the config. And will create a new credential under the path you just claimed
and set it to config, if it's the first time.

And also, you could specify other credentials by **configure** command.`,
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{"organization", "account", "environment", "credential"},
	RunE:      loginFunc,
}

func init() {
	command.Root.AddCommand(login)

	flagCred := constant.FlagCredential.ValStr()
	login.Flags().StringP(
		flagCred,
		"c",
		viper.GetString(flagCred),
		`
The credential file for the command about to be bound.
If it's the first time, or ignored, the default path will be used.'`)

	flagEnv := constant.FlagEnvironment.ValStr()
	login.Flags().StringP(
		flagEnv,
		"e",
		viper.GetString(flagEnv),
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
		return errors.New(fmt.Sprintf("reading passowrd error: %s\n", err.Error()))
	}

	if len(pwdBytes) == 0 {
		return errors.New("empty cryptPwd error")
	}

	success, err := supplierLogin(pwdBytes)
	if err != nil {
		return err
	}

	return syncLogin(success)
}

func supplierLogin(cryptPwdBytes []byte) (*resty.Response, error) {
	response, err := restyx.Client.R().
		EnableTrace().
		SetBody(ReqLogin{
			Account:      viper.GetString(constant.FlagAccount.ValStr()),
			Organization: viper.GetString(constant.FlagOrganization.ValStr()),
			Password:     cryptPwdBytes,
			Environment:  viper.GetString(constant.FlagEnvironment.ValStr()),
		}).
		Post(viper.GetString("bound_with.endpoint") + "/oauth/login")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("resty error: %s\n", err.Error()))
	}

	slog.Debug(fmt.Sprintf("response: %v\n", response)) // TODO: remove

	if e := util.HasError(response); e != nil {
		return nil, errors.New(fmt.Sprintf("supplier error: %s\n", e))
	}

	return response, nil
}

func syncLogin(response *resty.Response) error {
	cred := new(config.Credential)
	if err := json.Unmarshal(response.Body(), cred); err != nil {
		return errors.New(fmt.Sprintf("unmarshaling json response error: %s\n", err.Error()))
	}

	if err := config.WriteCredential(viper.GetString(constant.FlagCredential.ValStr()), cred); err != nil {
		return err
	}

	return config.SyncConfigByFlags()
}
