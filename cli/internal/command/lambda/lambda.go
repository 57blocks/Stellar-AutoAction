package lambda

import (
	"github.com/57blocks/auto-action/cli/internal/command"

	"github.com/spf13/cobra"
)

// lambdaGroup represents the lambda command
var lambdaGroup = &cobra.Command{
	Use:   "lambda",
	Short: "Lambda action management",
}

func init() {
	command.Root.AddCommand(lambdaGroup)
}
