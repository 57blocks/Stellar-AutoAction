package auth

import (
	"fmt"

	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/pkg/errorx"
	"github.com/57blocks/auto-action/cli/internal/pkg/logx"
	"github.com/57blocks/auto-action/cli/internal/pkg/restyx"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

// logout represents the logout command
var logout = &cobra.Command{
	Use:   "logout",
	Short: "End the current authenticated session",
	Long: `
Description:
  The logout command terminates the current authenticated session using the credential 
  path specified in the configuration file.

Behavior:
  - This command will invalidate the current session token.
  - The credential file associated with the current session will be cleared.

Notes:
  - This action only affects the currently active credential.
  - Other stored credentials remain valid and can be accessed using the 'configure' command.
  - To switch to a different authenticated session, use the 'configure' command to select 
    another credential instead of logging out and back in.

Examples:
  autoaction auth logout

Related Commands:
  autoaction auth configure - Manage multiple credentials and switch between sessions
`,
	Args: cobra.NoArgs,
	RunE: logoutFunc,
}

func init() {
	authGroup.AddCommand(logout)
}

type ReqLogout struct {
	Token string `json:"token"`
}

func logoutFunc(_ *cobra.Command, _ []string) error {
	cfg, err := config.ReadConfig()
	if err != nil {
		return err
	}

	if _, err := supplierLogout(); err != nil {
		return errorx.Internal(fmt.Sprintf("logging out error: %s", err.Error()))
	}

	if err := config.RemoveCredential(cfg.Credential); err != nil {
		return errorx.Internal(fmt.Sprintf("cleaning up error: %s", err.Error()))
	}

	if err := config.ResetConfigCredential(); err != nil {
		return errorx.Internal(fmt.Sprintf("resetting config error: %s", err.Error()))
	}

	logx.Logger.Info("you've logged out")

	return nil
}

func supplierLogout() (*resty.Response, error) {
	token, err := config.Token()
	if err != nil {
		logx.Logger.Error("PS: Should login first.")
		return nil, err
	}

	URL := util.ParseReqPath(fmt.Sprintf("%s/oauth/logout", config.Vp.GetString("bound_with.endpoint")))

	response, err := restyx.Client.R().
		EnableTrace().
		SetHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": token,
		}).
		Delete(URL)
	if err != nil {
		return nil, errorx.RestyError(err.Error())
	}
	if response.IsError() {
		return nil, errorx.WithRestyResp(response)
	}

	return response, nil
}
