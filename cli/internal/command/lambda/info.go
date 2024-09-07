package lambda

import (
	"fmt"

	"github.com/57blocks/auto-action/cli/internal/command"
	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/pkg/errorx"
	"github.com/57blocks/auto-action/cli/internal/pkg/logx"
	"github.com/57blocks/auto-action/cli/internal/pkg/restyx"

	"github.com/spf13/cobra"
)

// info represents the info command
var info = &cobra.Command{
	Use:   "info <name/arn>",
	Short: "Lambda essential information",
	Long: `
Description:
  Query the essential information of a specific Lambda, by name/arn.
  Which includes the VPC and Event Bridge Schedulers bound with.

Note:
  - The results contains the essential info about VPC and Schedulers.
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return infoFunc(cmd, args)
	},
}

func init() {
	command.Root.AddCommand(info)
}

func infoFunc(_ *cobra.Command, args []string) error {
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
		Get(fmt.Sprintf("%s/lambda/%s/info", config.Vp.GetString("bound_with.endpoint"), args[0]))
	if err != nil {
		return errorx.RestyError(err.Error())
	}
	if response.IsError() {
		return errorx.WithRestyResp(response)
	}

	logx.Logger.Info("lambda info", "result", response.String())

	return nil
}
