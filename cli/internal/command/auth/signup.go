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
	Short: "Create a new account for Stellar AutoAction",
	Long: `
Description:
  The signup command allows you to create a new account for the Stellar AutoAction system.
  You will need to provide your organization name, desired username, and a password.

Required Information:
  - Organization name
  - Username
  - Password

Process:
  1. Enter the required information when prompted.
  2. The system will validate your input and create your account.
  3. Upon successful signup, you can use the 'login' command to authenticate.

Notes:
  - The organization name must already exist in the system. An error will occur if it doesn't.
  - Usernames must be unique within an organization. Duplicate usernames are not allowed.
  - Ensure your password meets the system's security requirements.

Examples:
  autoaction auth signup -o "MyOrg" -a "john.doe" -d "Developer account"

Related Commands:
  autoaction auth login - Authenticate with your new account after signup
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
		`Username for the new account.
Must be unique within the organization.
Required for signup.`)

	flagOrg := constant.FlagOrganization.ValStr()
	signup.Flags().StringP(flagOrg,
		"o",
		"",
		`Name of the organization to join.
Must be an existing organization in the system.
Required for signup.`)

	flagDesc := constant.FlagDescription.ValStr()
	signup.Flags().StringP(flagDesc,
		"d",
		"",
		`Optional description for the user account.
Can be used to provide additional information about the user or their role.`)

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
