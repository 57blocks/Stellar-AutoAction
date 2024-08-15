package oauth

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"syscall"

	"github.com/57blocks/auto-action/cli/internal/command"
	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/constant"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"
	"github.com/57blocks/auto-action/cli/internal/third-party/req"

	"github.com/BurntSushi/toml"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
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
		"the credential file for the command about to be executed")

	flagEnv := constant.FlagEnvironment.ValStr()
	login.Flags().StringP(
		flagEnv,
		"e",
		viper.GetString(flagEnv),
		"the execution environment")

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

type (
	Req struct {
		Account      string `json:"account"`
		Organization string `json:"organization"`
		Password     []byte `json:"password"`
		Environment  string `json:"environment"`
	}
	Resp struct {
		Token   string `json:"token"`
		Refresh string `json:"refresh"`
	}
)

func loginFunc(cmd *cobra.Command, args []string) error {
	fmt.Println("Password: ")

	pwdBytes, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return errors.New(fmt.Sprintf("reading passowrd error: %s\n", err.Error()))
	}

	if len(pwdBytes) == 0 {
		return errors.New("empty cryptPwd error")
	}

	cryptPwdBytes, err := bcrypt.GenerateFromPassword(pwdBytes, bcrypt.DefaultCost)
	if err != nil {
		return errors.New(fmt.Sprintf("crypting pwd error: %s\n", err.Error()))
	}

	response, err := req.Client.R().
		EnableTrace().
		SetBody(Req{
			Account:      viper.GetString(constant.FlagAccount.ValStr()),
			Organization: viper.GetString(constant.FlagOrganization.ValStr()),
			Password:     cryptPwdBytes,
			Environment:  viper.GetString(constant.FlagEnvironment.ValStr()),
		}).
		Post(constant.Host.String() + "/oauth/login")
	if err != nil {
		return errors.New(fmt.Sprintf("endpoint request error: %s\n", err.Error()))
	}

	slog.Debug(fmt.Sprintf("response: %v\n", response))

	if e := util.IsError(response); e != nil {
		fmt.Println(e)
		return errors.New(fmt.Sprintf("response error: %s\n", e))
	}

	return syncCred(response)
}

func syncCred(response *resty.Response) error {
	cred := new(config.Credential)
	if err := json.Unmarshal(response.Body(), cred); err != nil {
		return errors.New(fmt.Sprintf("unmarshaling json response error: %s\n", err.Error()))
	}

	credToml, err := toml.Marshal(cred)
	if err != nil {
		return errors.New(fmt.Sprintf("marshaling credentials error: %s\n", err.Error()))
	}

	err = os.WriteFile(viper.GetString(constant.FlagCredential.ValStr()), credToml, 0666)
	if err != nil {
		return errors.New(fmt.Sprintf("writing credentials error: %s\n", err.Error()))
	}

	return config.SyncConfig(viper.ConfigFileUsed())
}
