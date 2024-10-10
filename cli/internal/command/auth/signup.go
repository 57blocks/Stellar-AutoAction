package auth

import (
	"fmt"
	"os"

	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/constant"
	"github.com/57blocks/auto-action/cli/internal/pkg/errorx"
	"github.com/57blocks/auto-action/cli/internal/pkg/logx"
	"github.com/57blocks/auto-action/cli/internal/pkg/restyx"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var signup = &cobra.Command{
	Use:   "signup",
	Short: "Signup to the Stellar auto-action.",
	Long: `
Description:
  To signup for the Stellar Auto-Action system, 
  you need to provide the organization name, username, and password during signup. 
  After successful signup, you can log in using the login command.

Note:
  - An error will occur if the organization name does not exist.
  - Duplicate usernames are not allowed within the same organization.
`,
	Args: cobra.NoArgs,
	RunE: signupFunc,
}

func init() {
	authGroup.AddCommand(signup)

	flagAcc := constant.FlagAccount.ValStr()
	signup.Flags().StringP(
		flagAcc,
		"a",
		"",
		"name of the account")

	flagOrg := constant.FlagOrganization.ValStr()
	signup.Flags().StringP(flagOrg,
		"o",
		"",
		"name of the organization")

	flagDesc := constant.FlagDescription.ValStr()
	signup.Flags().StringP(flagDesc,
		"d",
		"",
		"description of the user")

	if err := signup.MarkFlagRequired(flagAcc); err != nil {
		return
	}
	if err := signup.MarkFlagRequired(flagOrg); err != nil {
		return
	}
}

type ReqSignup struct {
	Account      string  `json:"account"`
	Organization string  `json:"organization"`
	Description  *string `json:"description,omitempty"`
	Password     string  `json:"password"`
}

func signupFunc(cmd *cobra.Command, args []string) error {
	fmt.Println("Password: ")

	pwdBytes, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return errorx.Internal(fmt.Sprintf("reading passowrd error: %s", err.Error()))
	}

	if len(pwdBytes) == 0 {
		return errorx.BadRequest("empty password error")
	}

	fmt.Println("Confirm password: ")
	confirmPwdBytes, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return errorx.Internal(fmt.Sprintf("reading confirm password error: %s", err.Error()))
	}

	if len(confirmPwdBytes) == 0 {
		return errorx.BadRequest("empty confirm password error")
	}

	if string(pwdBytes) != string(confirmPwdBytes) {
		return errorx.BadRequest("password and confirm password not match")
	}

	key, err := util.LoadPublicKey(config.Vp.GetString("general.public_key"))
	if err != nil {
		return err
	}
	encodedPwd, err := util.EncryptPassword(string(pwdBytes), key)
	if err != nil {
		return err
	}

	err = supplierSignup(encodedPwd)
	if err != nil {
		return err
	}

	logx.Logger.Info("Signup success! Please login. ")

	return nil
}

func supplierSignup(pwdHash string) error {
	URL := util.ParseReqPath(fmt.Sprintf("%s/oauth/signup", config.Vp.GetString("bound_with.endpoint")))

	description := config.Vp.GetString(constant.FlagDescription.ValStr())
	var descPtr *string
	if description != "" {
		descPtr = &description
	}

	response, err := restyx.Client.R().
		EnableTrace().
		SetBody(ReqSignup{
			Account:      config.Vp.GetString(constant.FlagAccount.ValStr()),
			Organization: config.Vp.GetString(constant.FlagOrganization.ValStr()),
			Description:  descPtr,
			Password:     pwdHash,
		}).
		Post(URL)
	if err != nil {
		return errorx.RestyError(err.Error())
	}
	if response.IsError() {
		return errorx.WithRestyResp(response)
	}

	return nil
}
