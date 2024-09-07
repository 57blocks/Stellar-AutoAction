package oauth

import (
	"fmt"
	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/pkg/errorx"
	"github.com/57blocks/auto-action/cli/internal/pkg/restyx"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"
	"log/slog"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

// logout represents the logout command
var logout = &cobra.Command{
	Use:   "logout",
	Short: "Logout the current session",
	Long: `
Description:
  Logout the current session by the path of credential in the config.

Note:
  - For other credentials, they are still alive. It's recommended that
    switching session by using **configure** command to set other 
    credentials.
`,
	Args: cobra.NoArgs,
	RunE: logoutFunc,
}

func init() {
	oauthGroup.AddCommand(logout)
}

type ReqLogout struct {
	Token string `json:"token"`
}

func logoutFunc(_ *cobra.Command, _ []string) error {
	cfg, err := config.ReadConfig()
	if err != nil {
		return err
	}

	// logout already
	if cfg.Credential == "" {
		slog.Info("you've already logged out")
		return nil
	}

	// credential does not exist
	if !util.IsExists(cfg.Credential) {
		slog.Info("credential not found, reset the config directly.")
		return config.ResetConfigCredential()
	}

	// logout
	credential, err := config.ReadCredential(cfg.Credential)
	if err != nil {
		return err
	}

	if _, err := supplierLogout(credential.Token); err != nil {
		return err
	}

	if err := config.RemoveCredential(cfg.Credential); err != nil {
		return err
	}

	if err := config.ResetConfigCredential(); err != nil {
		return err
	}

	return nil
}

func supplierLogout(token string) (*resty.Response, error) {
	URL := util.ParseReqPath(fmt.Sprintf("%s/oauth/logout", config.Vp.GetString("bound_with.endpoint")))

	response, err := restyx.Client.R().
		EnableTrace().
		SetBody(ReqLogout{Token: token}).
		Delete(URL)
	if err != nil {
		return nil, errorx.WithRestyResp(response)
	}
	if response.IsError() {
		return nil, errorx.WithRestyResp(response)
	}

	return response, nil
}
