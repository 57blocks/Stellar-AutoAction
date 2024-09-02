package lambda

import (
	"fmt"
	"log/slog"

	"github.com/57blocks/auto-action/cli/internal/command"
	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/pkg/restyx"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// info represents the info command
var info = &cobra.Command{
	Use:   "info <name/arn>",
	Short: "Lambda essential information",
	Long: `
Description:
  Query the essential information of a specific Lambda, by name/arn.
  Which includes the VPC and Event Bridge Schedulers bound with.
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
		Get(fmt.Sprintf("%s/lambda/%s/info", viper.GetString("bound_with.endpoint"), args[0]))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("resty error: %s\n", err.Error()))
	}
	if e := util.HasError(response); e != nil {
		return errors.Wrap(e, fmt.Sprintf("supplier error: %s\n", e))
	}

	slog.Debug(fmt.Sprintf("%v\n", response))

	return nil
}
