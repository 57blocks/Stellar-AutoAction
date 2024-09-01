package lambda

import (
	"fmt"

	"github.com/57blocks/auto-action/cli/internal/command"
	"github.com/57blocks/auto-action/cli/internal/constant"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// trigger represents the trigger command
var trigger = &cobra.Command{
	Use:   "trigger <name/arn>",
	Short: "Trigger a specific lambda function",
	Long: `Trigger a specific lambda function by its ARN, which is inputted as an
argument in CLI. Then the lambda will be executed instantly.

As for the handlers arguments, `,
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
		"",
		viper.GetString(flagPayload),
		`The input event payload of the Lambda function. It should be a 
well-formed JSON string.

Meanwhile, it should be suitable/executable/valid in your handler to use.'`)
}

func triggerFunc(_ *cobra.Command, args []string) error {
	fmt.Println("trigger called")

	return nil
}
