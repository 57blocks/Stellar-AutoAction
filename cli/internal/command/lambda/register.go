package lambda

import (
	"fmt"
	"strings"

	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/constant"
	"github.com/57blocks/auto-action/cli/internal/pkg/errorx"
	"github.com/57blocks/auto-action/cli/internal/pkg/logx"
	"github.com/57blocks/auto-action/cli/internal/pkg/restyx"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
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
			return errorx.BadRequest("at most one expression flag should be set")
		}

		return nil
	},
	Args: cobra.MinimumNArgs(1),
	RunE: registerFunc,
}

func init() {
	lambdaGroup.AddCommand(register)

	flagCron := constant.FlagCron.ValStr()
	register.Flags().StringP(
		flagCron,
		"c",
		config.Vp.GetString(flagCron),
		`The cron execution expression for the Event Bridge Scheduler.
For more info: https://docs.aws.amazon.com/scheduler/latest/UserGuide/schedule-types.html#rate-based
`)

	flagRate := constant.FlagRate.ValStr()
	register.Flags().StringP(
		flagRate,
		"r",
		config.Vp.GetString(flagRate),
		`The rate execution expression for the Event Bridge Scheduler.
For more info: https://docs.aws.amazon.com/scheduler/latest/UserGuide/schedule-types.html#cron-based
`)

	flagAt := constant.FlagAt.ValStr()
	register.Flags().StringP(
		flagAt,
		"a",
		config.Vp.GetString(flagAt),
		`The one-time execution expression for the Event Bridge Scheduler.
For more info: https://docs.aws.amazon.com/scheduler/latest/UserGuide/schedule-types.html#one-time
`)
}

func registerFunc(_ *cobra.Command, args []string) error {
	resp, err := supplierRegister(args)
	if err != nil {
		return err
	}

	logx.Logger.Info("register lambda", "result", resp.String())

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

	URL := util.ParseReqPath(fmt.Sprintf("%s/lambda", config.Vp.GetString("bound_with.endpoint")))

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
		if value := strings.TrimSpace(config.Vp.GetString(key)); value != "" {
			flagMap["expression"] = value
			logx.Logger.Info("register lambda", "invoke expression", value)
			break
		}
	}

	if len(flagMap) == 0 {
		flagMap["expression"] = ""
		logx.Logger.Info("register lambda", "invoke expression", "manually")
	}

	request = request.SetFormData(flagMap)

	response, err := request.Post(URL)
	if err != nil {
		return nil, errorx.WithRestyResp(response)
	}
	if response.IsError() {
		return nil, errorx.WithRestyResp(response)
	}

	return response, nil
}
