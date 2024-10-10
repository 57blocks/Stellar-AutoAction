package action

import (
	"encoding/json"
	"fmt"

	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/pkg/errorx"
	"github.com/57blocks/auto-action/cli/internal/pkg/logx"
	"github.com/57blocks/auto-action/cli/internal/pkg/restyx"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"

	"github.com/spf13/cobra"
)

// info represents the info command
var info = &cobra.Command{
	Use:   "info <name/arn>",
	Short: "Display essential information about an action",
	Long: `
Description:
  The info command retrieves and displays essential information about a specific action,
  identified by its name or ARN (Amazon Resource Name).

This command provides details including:
  - Basic action configuration
  - Bound Event Bridge Schedulers

Arguments:
  <name/arn>    The name or ARN of the action to query

Examples:
  autoaction action info my-action
  autoaction action info arn:aws:lambda:us-west-2:123456789012:function:my-action

Note:
  - The output includes key information about the associated Schedulers.
  - Ensure you have the necessary permissions to view the action details.
`,
	Args: cobra.ExactArgs(1),
	RunE: infoFunc,
}

func init() {
	actionGroup.AddCommand(info)
}

func infoFunc(_ *cobra.Command, args []string) error {
	token, err := config.Token()
	if err != nil {
		logx.Logger.Error("PS: Should login first.")
		return err
	}

	URL := util.ParseReqPath(fmt.Sprintf("%s/lambda/%s", config.Vp.GetString("bound_with.endpoint"), args[0]))

	response, err := restyx.Client.R().
		EnableTrace().
		SetHeaders(map[string]string{
			"Content-Type":  "multipart/form-data",
			"Authorization": token,
		}).
		Get(URL)
	if err != nil {
		return errorx.RestyError(err.Error())
	}
	if response.IsError() {
		return errorx.WithRestyResp(response)
	}

	var respData map[string]interface{}
	if err := json.Unmarshal(response.Body(), &respData); err != nil {
		logx.Logger.Error("Error unmarshalling JSON", "error", err.Error())
		return errorx.Internal(err.Error())
	}

	logx.Logger.Info("action info", "result", respData)

	return nil
}
