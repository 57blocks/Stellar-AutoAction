package lambda

import (
	"fmt"
	"log/slog"

	"github.com/57blocks/auto-action/cli/internal/command"
	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/constant"
	"github.com/57blocks/auto-action/cli/internal/pkg/restyx"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
  1. If the Lambda does not depend on the input in the EVENT, the 
  payload is not required.
  2. If so, the payload should be a well-formed JSON string, which is
  suitable/executable/valid in your handler to use.
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
		viper.GetString(flagPayload),
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
			"payload": viper.GetString(constant.FlagPayload.ValStr()),
		}).
		Post(fmt.Sprintf("%s/lambda/%s", viper.GetString("bound_with.endpoint"), args[0]))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("resty error: %s\n", err.Error()))
	}
	if e := util.HasError(response); e != nil {
		return errors.Wrap(e, fmt.Sprintf("supplier error: %s\n", e))
	}

	slog.Debug(fmt.Sprintf("%v\n", response))

	return nil
}
