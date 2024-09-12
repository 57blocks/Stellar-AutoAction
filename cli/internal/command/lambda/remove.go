package lambda

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

// removeCmd represents the lambda remove command
var removeCmd = &cobra.Command{
	Use:   "remove <name/arn>",
	Short: "Remove a specific lambda function by its name/ARN",
	Long: `
Description:
	Remove a specific lambda function by its name/ARN, together with
its trigger if it has one.

Note: 
  - The execution logs are still on the CloudWatch Logs.
  - Each of the Lambda bound with one Scheduler to get triggered, so 
    the remove response should contains the info of them both.
`,
	Args: cobra.ExactArgs(1),
	RunE: removeFunc,
}

func init() {
	lambdaGroup.AddCommand(removeCmd)
}

func removeFunc(cmd *cobra.Command, args []string) error {
	token, err := config.Token()
	if err != nil {
		return err
	}

	URL := util.ParseReqPath(fmt.Sprintf("%s/lambda/%s", config.Vp.GetString("bound_with.endpoint"), args[0]))

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
