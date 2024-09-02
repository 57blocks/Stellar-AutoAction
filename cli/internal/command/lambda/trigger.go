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

// trigger represents the trigger command
var trigger = &cobra.Command{
	Use:   "trigger <name/arn> [flags]",
	Short: "Trigger a specific lambda function with its required payload",
	Long: `
Trigger a specific lambda function by its ARN, which is inputted as an
argument in CLI. Then the lambda will be executed instantly.

As for the handlers required payload, it should be a well-formed JSON 
string, suitable/executable/valid in your handler to use.
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return triggerFunc(cmd, args)
	},
}

func init() {
	command.Root.AddCommand(trigger)

	flagPayload := constant.FlagPayload.ValStr()
	trigger.Flags().StringP(
		flagPayload,
		"p",
		viper.GetString(flagPayload),
		`The input event payload of the Lambda function. It should be a 
well-formed JSON string.

Meanwhile, it should be suitable/executable/valid in your handler to use.'`)
}

func triggerFunc(_ *cobra.Command, args []string) error {
	token, err := config.Token()
	if err != nil {
		return err
	}

	response, err := restyx.Client.R().
		EnableTrace().
		SetHeaders(map[string]string{
			"Content-Type":  "multipart/form-data",
			"Authorization": token,
		}).
		Post(fmt.Sprintf("%s/lambda/%s", viper.GetString("bound_with.endpoint"), args[0]))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("resty error: %s\n", err.Error()))
	}

	slog.Debug(fmt.Sprintf("%v\n", response)) // TODO: remove

	if e := util.HasError(response); e != nil {
		return errors.Wrap(e, fmt.Sprintf("supplier error: %s\n", e))
	}

	return nil
}
