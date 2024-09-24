package lambda

import (
	"encoding/json"
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
  - Payload, using -p/--payload to set the payload, if it's required.
    The payload is a well-formed JSON string, and also, be valid/usable
    in the handler/handlers.
  - Trigger:
	- By cron.
	- By rate, only three units supported: minutes, hours, days.
      For example: rate(1 minutes).
	- By at, one-time execution, under a specific time in the future.
      For example: at(yyyy-mm-ddThh:mm:ss).
	- At most one expression flag could be set.
	- Expression flags: at/cron/rate would create an Event Bridge
      Scheduler to invoke lambda function, together with the payload.
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

	fCron := constant.FlagCron.ValStr()
	register.Flags().StringP(
		fCron,
		"c",
		config.Vp.GetString(fCron),
		`The cron execution expression for the Event Bridge Scheduler.
For more info: https://docs.aws.amazon.com/scheduler/latest/UserGuide/schedule-types.html#cron-based
`)

	fRate := constant.FlagRate.ValStr()
	register.Flags().StringP(
		fRate,
		"r",
		config.Vp.GetString(fRate),
		`The rate execution expression for the Event Bridge Scheduler.
For more info: https://docs.aws.amazon.com/scheduler/latest/UserGuide/schedule-types.html#rate-based
`)

	fAt := constant.FlagAt.ValStr()
	register.Flags().StringP(
		fAt,
		"a",
		config.Vp.GetString(fAt),
		`The one-time execution expression for the Event Bridge Scheduler.
For more info: https://docs.aws.amazon.com/scheduler/latest/UserGuide/schedule-types.html#one-time
`)

	fPayload := constant.FlagPayload.ValStr()
	register.Flags().StringP(
		fPayload,
		"p",
		config.Vp.GetString(fPayload),
		`The payload for the lambda function execution.
`)
}

func registerFunc(_ *cobra.Command, args []string) error {
	if err := util.ValidateZipFiles(args); err != nil {
		return err
	}

	response, err := supplierRegister(args)
	if err != nil {
		return err
	}

	var respData []map[string]interface{}
	if err := json.Unmarshal(response.Body(), &respData); err != nil {
		logx.Logger.Error("Error unmarshalling JSON", "error", err.Error())
		return errorx.Internal(err.Error())
	}

	logx.Logger.Info("lambda info", "result", respData)

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
	fMap := make(map[string]string, 1)
	fKeys := []string{constant.FlagAt.ValStr(), constant.FlagCron.ValStr(), constant.FlagRate.ValStr()}
	for _, key := range fKeys {
		if expVal := strings.TrimSpace(config.Vp.GetString(key)); expVal != "" {
			fMap["expression"] = expVal
			logx.Logger.Info("register lambda", "invoke expression", expVal)
			break
		}
	}

	if len(fMap) == 0 {
		logx.Logger.Info("register lambda", "invoke expression", "manually")
	}

	if pldVal := strings.TrimSpace(config.Vp.GetString(constant.FlagPayload.ValStr())); pldVal != "" {
		temp := make(map[string]string, 1)
		if err := json.Unmarshal([]byte(pldVal), &temp); err != nil {
			return nil, errorx.BadRequest(fmt.Sprintf("invalid payload: %s", pldVal))
		}

		fMap["payload"] = pldVal
	}
	logx.Logger.Info("register lambda", "invoke payload", "none")

	request = request.SetFormData(fMap)

	response, err := request.Post(URL)
	if err != nil {
		return nil, errorx.RestyError(err.Error())
	}
	if response.IsError() {
		return nil, errorx.WithRestyResp(response)
	}

	return response, nil
}
