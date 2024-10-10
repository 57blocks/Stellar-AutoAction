package action

import (
	"encoding/json"
	"fmt"

	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/constant"
	"github.com/57blocks/auto-action/cli/internal/pkg/errorx"
	"github.com/57blocks/auto-action/cli/internal/pkg/logx"
	"github.com/57blocks/auto-action/cli/internal/pkg/restyx"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"

	"github.com/spf13/cobra"
)

// invoke represents the invoke command
var invoke = &cobra.Command{
	Use:   "invoke <name/arn> [flags]",
	Short: "Execute a specific action with an optional payload",
	Long: `
Description:
  The invoke command allows you to execute a specific action immediately by providing
  its name or ARN (Amazon Resource Name). This command is useful for testing actions
  or triggering them manually.

Arguments:
  <name/arn>    The name or ARN of the action to invoke

Examples:
  autoaction action invoke my-action
  autoaction action invoke arn:aws:lambda:us-west-2:123456789012:function:my-action
  autoaction action invoke my-action -p '{"key": "value"}'

Notes:
  - If the action does not require input data, the payload flag can be omitted.
  - When provided, the payload must be a valid JSON string that matches the
    expected input format of your action's handler.
  - Ensure you have the necessary permissions to invoke the action.
`,
	Args: cobra.ExactArgs(1),
	RunE: invokeFunc,
}

func init() {
	actionGroup.AddCommand(invoke)

	flagPayload := constant.FlagPayload.ValStr()
	invoke.Flags().StringP(
		flagPayload,
		"p",
		config.Vp.GetString(flagPayload),
		`A well-formed JSON string representing the payload for the action. 
This payload should be compatible with your action's handler. 
Example: '{"key": "value"}'
`)
}

func invokeFunc(_ *cobra.Command, args []string) error {
	token, err := config.Token()
	if err != nil {
		logx.Logger.Error("PS: Should login first.")
		return err
	}

	URL := util.ParseReqPath(fmt.Sprintf("%s/lambda/%s", config.Vp.GetString("bound_with.endpoint"), args[0]))

	response, err := restyx.Client.R().
		EnableTrace().
		SetHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": token,
		}).
		SetBody(map[string]string{
			"payload": config.Vp.GetString(constant.FlagPayload.ValStr()),
		}).
		Post(URL)
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

	logx.Logger.Info("invoke action success", "result", respData)

	return nil
}
