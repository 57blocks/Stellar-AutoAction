package lambda

import (
	"fmt"

	"github.com/57blocks/auto-action/cli/internal/command"
	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/constant"
	"github.com/57blocks/auto-action/cli/internal/pkg/errorx"
	"github.com/57blocks/auto-action/cli/internal/pkg/logx"
	"github.com/57blocks/auto-action/cli/internal/pkg/restyx"

	"github.com/spf13/cobra"
)

// invoke represents the invoke command
var invoke = &cobra.Command{
	Use:   "invoke <name/arn> [flags]",
	Short: "Invoke a specific lambda function with its required payload",
	Long: `
Description:
  Invoke a specific lambda function by its name/ARN, which is inputted
  as an argument. Then the lambda will be executed instantly.

Note:
  - If the Lambda does not depend on the input in the EVENT, the 
    payload is not required.
  - If so, the payload should be a well-formed JSON string, which is
    suitable/executable/valid in your handler event to use.
    For example: -p '{"key": "value"}'
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return invokeFunc(cmd, args)
	},
}

func init() {
	command.Root.AddCommand(invoke)

	flagPayload := constant.FlagPayload.ValStr()
	invoke.Flags().StringP(
		flagPayload,
		"p",
		config.Vp.GetString(flagPayload),
		`
A well-formed JSON string. And should be suitable/executable/valid in
your handler to use. Example: '{"key": "value"}'
`)
}

func invokeFunc(_ *cobra.Command, args []string) error {
	token, err := config.Token()
	if err != nil {
		return err
	}

	response, err := restyx.Client.R().
		EnableTrace().
		SetHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": token,
		}).
		SetBody(map[string]string{
			"payload": config.Vp.GetString(constant.FlagPayload.ValStr()),
		}).
		Post(fmt.Sprintf("%s/lambda/%s", config.Vp.GetString("bound_with.endpoint"), args[0]))
	if err != nil {
		return errorx.RestyError(err.Error())
	}
	if response.IsError() {
		return errorx.WithRestyResp(response)
	}

	logx.Logger.Info("invoke lambda success", "result", response.String())

	return nil
}
