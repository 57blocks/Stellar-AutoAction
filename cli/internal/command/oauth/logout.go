package oauth

import (
	"encoding/json"
	"fmt"
	"github.com/57blocks/auto-action/cli/internal/credential"
	"log/slog"
	"os"

	"github.com/57blocks/auto-action/cli/internal/command"
	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/constant"
	"github.com/57blocks/auto-action/cli/internal/pkg/restyx"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"

	"github.com/BurntSushi/toml"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// logout represents the logout command
var logout = &cobra.Command{
	Use:   "logout",
	Short: "Logout the current session",
	Long: `Logout the current session under the path of credential in
the config.

For other credentials, they are still alive. It's recommended that
switching session by using **configure** command to set other 
credentials to the config`,
	Args: cobra.NoArgs,
	RunE: logoutFunc,
}

func init() {
	command.Root.AddCommand(logout)
}

type ReqLogout struct {
	Token string `json:"token"`
}

func logoutFunc(_ *cobra.Command, _ []string) error {
	cfg, err := config.ReadConfig(util.DefaultPath())
	if err != nil {
		return errors.New(fmt.Sprintf("reading configuration error: %s\n", err.Error()))
	}

	if util.IsExists(cfg.Credential) {
		credential, err := credential.ReadCredential(cfg.Credential)
		if err != nil {
			return errors.New(fmt.Sprintf("reading credential error: %s\n", err.Error()))
		}

		response, err := supplierLogout(credential.Token)
		if err != nil {
			return errors.New(fmt.Sprintf("supplier error: %s\n", err.Error()))
		}
		if e := util.HasError(response); e != nil {
			return errors.New(fmt.Sprintf("supplier error: %s\n", e))
		}

		if err := os.Remove(cfg.Credential); err != nil {
			return errors.New(fmt.Sprintf("removing credential error: %s\n", err.Error()))
		}
	} else {
		slog.Info("credential file not found")
	}

	cfg.Credential = ""
	if err := config.WriteConfig(cfg, viper.ConfigFileUsed()); err != nil {
		return errors.New(fmt.Sprintf("writing configuration error: %s\n", err.Error()))
	}

	slog.Info("logout successfully")

	return nil
}

func supplierLogout(token string) (*resty.Response, error) {
	response, err := restyx.Client.R().
		EnableTrace().
		SetBody(ReqLogout{Token: token}).
		Post(viper.GetString("bound.endpoint") + "/oauth/logout")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("resty request error: %s\n", err.Error()))
	}

	slog.Debug(fmt.Sprintf("response: %v\n", response)) // TODO: remove

	if e := util.HasError(response); e != nil {
		return nil, errors.New(fmt.Sprintf("supplier error: %s\n", e))
	}

	return response, nil
}

func syncLogout(response *resty.Response) error {
	cred := new(credential.Credential)
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

	return config.SyncConfigByFlags(viper.ConfigFileUsed())
}
