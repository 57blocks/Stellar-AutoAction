package oauth

import (
	"fmt"
	"log/slog"

	"github.com/57blocks/auto-action/cli/internal/command"
	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/pkg/restyx"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		return logoutFunc(cmd, args)
	},
}

func init() {
	command.Root.AddCommand(logout)
}

type ReqLogout struct {
	Token string `json:"token"`
}

func logoutFunc(_ *cobra.Command, _ []string) error {
	cfg, err := config.ReadConfig()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("reading config error: %s\n", err.Error()))
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
		return errors.New(fmt.Sprintf("reading credential error: %s\n", err.Error()))
	}

	if _, err := supplierLogout(credential.Token); err != nil {
		return errors.New(fmt.Sprintf("resty error: %s\n", err.Error()))
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
	response, err := restyx.Client.R().
		EnableTrace().
		SetBody(ReqLogout{Token: token}).
		Delete(viper.GetString("bound_with.endpoint") + "/oauth/logout")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("resty error: %s\n", err.Error()))
	}

	slog.Debug(fmt.Sprintf("response: %v\n", response)) // TODO: remove

	if e := util.HasError(response); e != nil {
		return nil, errors.New(fmt.Sprintf("supplier error: %s\n", e))
	}

	return response, nil
}
