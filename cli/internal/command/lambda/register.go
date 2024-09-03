package lambda

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/57blocks/auto-action/cli/internal/command"
	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/constant"
	"github.com/57blocks/auto-action/cli/internal/pkg/restyx"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// register represents the register command
var register = &cobra.Command{
	Use:   "register [zips/packages]",
	Short: "Lambda register of local handler",
	Long: `
Description:
  Register the local handler/handlers to Amazon Lambda, with the
  recurring/scheduled rule. 
	
Note:
  - The name of the Lambda is based on the file name, make it unique.
  - The name of the handler must be: **handler**. 
  - Manually, if no flags puts in, which means the handler/handlers 
	will be triggered by invoke command manually.
  - By corn. 
  - By rate, only three units supported: minutes, hours, days.
  	For example: rate(1 minutes).
  - By at, only one-time execution, for a specific time in the future.
  	For example: at(yyyy-mm-ddThh:mm:ss).
  - At most one expression flag could be set.
  - Only cron/rate/at would create an Event Bridge Scheduler to invoke 
	the lambda function.
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		a := cmd.Flags().Changed(constant.FlagAt.ValStr())
		c := cmd.Flags().Changed(constant.FlagCron.ValStr())
		r := cmd.Flags().Changed(constant.FlagRate.ValStr())

		if a && r || a && c || r && c {
			return errors.New("at most one expression flag should be set")
		}

		return nil
	},
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return registerFunc(cmd, args)
	},
}

func init() {
	command.Root.AddCommand(register)

	flagCron := constant.FlagCron.ValStr()
	register.Flags().StringP(
		flagCron,
		"c",
		viper.GetString(flagCron),
		`The cron execution expression for the Event Bridge Scheduler.
For more info: https://docs.aws.amazon.com/scheduler/latest/UserGuide/schedule-types.html#rate-based
`)

	flagRate := constant.FlagRate.ValStr()
	register.Flags().StringP(
		flagRate,
		"r",
		viper.GetString(flagRate),
		`The rate execution expression for the Event Bridge Scheduler.
For more info: https://docs.aws.amazon.com/scheduler/latest/UserGuide/schedule-types.html#cron-based
`)

	flagAt := constant.FlagAt.ValStr()
	register.Flags().StringP(
		flagAt,
		"a",
		viper.GetString(flagAt),
		`The one-time execution expression for the Event Bridge Scheduler.
For more info: https://docs.aws.amazon.com/scheduler/latest/UserGuide/schedule-types.html#one-time
`)
}

func registerFunc(_ *cobra.Command, args []string) error {
	resp, err := supplierRegister(args)
	if err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("lambda identifiers: %v\n", resp.String()))

	return nil
}

func supplierRegister(args []string) (*resty.Response, error) {
	token, err := config.Token()
	if err != nil {
		return nil, err
	}

	argMap := make(map[string]string, len(args))
	for _, path := range args {
		argMap[path] = path
	}

	request := restyx.Client.R().
		EnableTrace().
		SetHeaders(map[string]string{
			"Content-Type":  "multipart/form-data",
			"Authorization": token,
		}).
		SetFiles(argMap)

	// flags handling
	flagMap := make(map[string]string, 1)
	flagKeys := []string{constant.FlagAt.ValStr(), constant.FlagCron.ValStr(), constant.FlagRate.ValStr()}
	for _, key := range flagKeys {
		if value := strings.TrimSpace(viper.GetString(key)); value != "" {
			flagMap["expression"] = value
			slog.Debug(fmt.Sprintf("%s expression: %s\n", key, value))
			break
		}
	}

	if len(flagMap) == 0 {
		flagMap["expression"] = ""
		slog.Debug("lambda will be triggered manually\n")
	}

	request = request.SetFormData(flagMap)

	resp, err := request.Post(viper.GetString("bound_with.endpoint") + "/lambda")
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("resty error: %s\n", err.Error()))
	}
	if e := util.HasError(resp); e != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("supplier error: %s\n", e))
	}

	return resp, nil
}
