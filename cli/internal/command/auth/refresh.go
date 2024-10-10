package auth

import (
	"encoding/json"
	"fmt"

	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/pkg/errorx"
	"github.com/57blocks/auto-action/cli/internal/pkg/logx"
	"github.com/57blocks/auto-action/cli/internal/pkg/restyx"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

// refreshCmd represents the refresh command
var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Renew the expired access token",
	Long: `
Description:
  The refresh command renews an expired access token using the refresh token 
  stored in the credential file.

Behavior:
  - This command uses the stored refresh token to obtain a new access token.
  - The refresh token itself is not updated or changed during this process.

Notes:
  - The expiration time of the refresh token remains unchanged after this operation.
  - If the refresh token has expired, you will need to perform a full login again 
    using the 'autoaction auth login' command.
  - Regular use of this command can help maintain continuous access without 
    frequent full logins.

Example:
  autoaction auth refresh

Related Commands:
  autoaction auth login - Perform a full login when the refresh token expires
`,
	Args: cobra.NoArgs,
	RunE: refreshFunc,
}

func init() {
	authGroup.AddCommand(refreshCmd)
}

type ReqRefresh struct {
	Refresh string `json:"refresh"`
}

func refreshFunc(_ *cobra.Command, _ []string) error {
	cfg, err := config.ReadConfig()
	if err != nil {
		return err
	}

	response, err := supplierRefresh()
	if err != nil {
		return err
	}

	return syncRefresh(cfg, response)
}

func supplierRefresh() (*resty.Response, error) {
	refresh, err := config.RefreshToken()
	if err != nil {
		logx.Logger.Error("PS: Should login first.")
		return nil, err
	}

	URL := util.ParseReqPath(fmt.Sprintf("%s/oauth/refresh", config.Vp.GetString("bound_with.endpoint")))

	response, err := restyx.Client.R().
		EnableTrace().
		SetHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": refresh,
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

func syncRefresh(cfg *config.GlobalConfig, resp *resty.Response) error {
	cred := new(config.Credential)
	if err := json.Unmarshal(resp.Body(), cred); err != nil {
		return errorx.Internal(fmt.Sprintf("unmarshaling json response error: %s", err.Error()))
	}

	if err := config.WriteCredential(cfg.Credential, cred); err != nil {
		return err
	}

	logx.Logger.Info("refreshed")

	return nil
}
