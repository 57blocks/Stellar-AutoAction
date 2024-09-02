package oauth

import (
	"encoding/json"
	"fmt"
	"github.com/57blocks/auto-action/cli/internal/command"
	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/pkg/restyx"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// refreshCmd represents the refresh command
var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh expired access token.",
	Long: `
Description:
  Using the refresh token which stores in the credential file to refresh.

Note:
  - Refresh do not update the refresh token, so the expiration of it
    is still the same.
  - When the refresh token is expired, you need to login again.
`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return refreshFunc(cmd, args)
	},
}

func init() {
	command.Root.AddCommand(refreshCmd)

}

type ReqRefresh struct {
	Refresh string `json:"refresh"`
}

func refreshFunc(_ *cobra.Command, _ []string) error {
	cfg, err := config.ReadConfig()
	if err != nil {
		return errors.New(fmt.Sprintf("reading config error: %s\n", err.Error()))
	}

	// logout check
	if cfg.Credential == "" {
		return errors.New("you've already logged out")
	}

	// credential does not exist
	if !util.IsExists(cfg.Credential) {
		return errors.New("credential file does not exist")
	}

	// refresh
	credential, err := config.ReadCredential(cfg.Credential)
	if err != nil {
		return errors.New(fmt.Sprintf("reading credential error: %s\n", err.Error()))
	}

	response, err := supplierRefresh(credential.Refresh)
	if err != nil {
		return errors.New(fmt.Sprintf("resty error: %s\n", err.Error()))
	}

	return syncRefresh(cfg, response)
}

func supplierRefresh(refresh string) (*resty.Response, error) {
	response, err := restyx.Client.R().
		EnableTrace().
		SetBody(ReqRefresh{Refresh: refresh}).
		Post(viper.GetString("bound_with.endpoint") + "/oauth/refresh")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("resty error: %s\n", err.Error()))
	}

	if e := util.HasError(response); e != nil {
		return nil, errors.New(fmt.Sprintf("supplier error: %s\n", e))
	}

	return response, nil
}

func syncRefresh(cfg *config.GlobalConfig, resp *resty.Response) error {
	cred := new(config.Credential)
	if err := json.Unmarshal(resp.Body(), cred); err != nil {
		return errors.New(fmt.Sprintf("unmarshaling json response error: %s\n", err.Error()))
	}

	if err := config.WriteCredential(cfg.Credential, cred); err != nil {
		return err
	}

	return nil
}
