package action

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/pkg/errorx"
	"github.com/57blocks/auto-action/cli/internal/pkg/logx"
	"github.com/57blocks/auto-action/cli/internal/pkg/restyx"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"

	"github.com/spf13/cobra"
)

// removeCmd represents the action remove command
var removeCmd = &cobra.Command{
	Use:   "remove <name/arn>",
	Short: "Remove a specific action and its associated trigger",
	Long: `
Description:
  The remove command deletes a specific action from Stellar AutoAction, identified by its name or ARN 
  (Amazon Resource Name). If the action has an associated trigger (scheduler), it will also be removed.

Arguments:
  <name/arn>    The name or ARN of the action to remove

Examples:
  autoaction action remove my-action
  autoaction action remove arn:aws:lambda:us-west-2:123456789012:function:my-action

Notes:
  - This command removes both the action and its associated EventBridge Scheduler (if any).
  - Execution logs for the action will remain in CloudWatch Logs and are not deleted by this command.
  - The removal response will include information about both the action and its associated scheduler.

Caution:
  This operation is irreversible. Make sure you want to permanently remove the action before proceeding.
`,
	Args: cobra.ExactArgs(1),
	RunE: removeFunc,
}

func init() {
	actionGroup.AddCommand(removeCmd)
}

func removeFunc(cmd *cobra.Command, args []string) error {
	token, err := config.Token()
	if err != nil {
		logx.Logger.Error("PS: Should login first.")
		return err
	}

	URL := util.ParseReqPath(fmt.Sprintf("%s/lambda/%s", config.Vp.GetString("bound_with.endpoint"), url.PathEscape(args[0])))

	response, err := restyx.Client.R().
		EnableTrace().
		SetHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": token,
		}).
		Delete(URL)
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

	logx.Logger.Info("removed successfully", "removed", respData)

	return nil
}
