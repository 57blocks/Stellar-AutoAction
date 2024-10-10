package action

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
	Short: "Register local handlers as actions",
	Long: `
Description:
  Register local handler(s) to AutoAction, optionally with recurring or scheduled execution rules.

Notes:
  - Action name is derived from the file name; ensure it's unique.
  - The handler function must be named "handler".
  - Without flags, the action will be triggered manually via the invoke command.
  - Payload must be a valid JSON string, usable by the handler(s).
  - Only one scheduling expression (cron/rate/at) can be set per action.

Scheduling Options:
  - Cron: Standard cron expression
  - Rate: Supports minutes, hours, days. E.g., rate(1 minutes)
  - At: One-time execution. Format: at(yyyy-mm-ddThh:mm:ss)

Examples:
  autoaction action register ./handler.zip
  autoaction action register ./handler.zip -a 'at(2022-12-31T23:59:59)' -p '{"key": "value"}'
  autoaction action register ./handler.zip -r 'rate(1 minutes)' -p '{"key": "value"}'
  autoaction action register ./handler.zip -c 'cron(0 12 * * ? *)' -p '{"key": "value"}'
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
	actionGroup.AddCommand(register)

	fCron := constant.FlagCron.ValStr()
	register.Flags().StringP(
		fCron,
		"c",
		config.Vp.GetString(fCron),
		`Cron expression for scheduled execution.
Format: cron(Minutes Hours Day-of-month Month Day-of-week Year)
Example: cron(0 12 * * ? *) for daily at 12:00 PM UTC
More info: https://docs.aws.amazon.com/scheduler/latest/UserGuide/schedule-types.html#cron-based
`)

	fRate := constant.FlagRate.ValStr()
	register.Flags().StringP(
		fRate,
		"r",
		config.Vp.GetString(fRate),
		`Rate expression for recurring execution.
Format: rate(value unit)
Supported units: minutes, hours, days
Example: rate(5 minutes) for every 5 minutes
More info: https://docs.aws.amazon.com/scheduler/latest/UserGuide/schedule-types.html#rate-based
`)

	fAt := constant.FlagAt.ValStr()
	register.Flags().StringP(
		fAt,
		"a",
		config.Vp.GetString(fAt),
		`One-time execution at a specific future time.
Format: at(yyyy-mm-ddThh:mm:ss)
Example: at(2023-12-31T23:59:59) for Dec 31, 2023 at 23:59:59 UTC
More info: https://docs.aws.amazon.com/scheduler/latest/UserGuide/schedule-types.html#one-time
`)

	fPayload := constant.FlagPayload.ValStr()
	register.Flags().StringP(
		fPayload,
		"p",
		config.Vp.GetString(fPayload),
		`JSON payload for the action execution.
Must be a valid JSON string.
Example: '{"key": "value"}'
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

	logx.Logger.Info("action info", "result", respData)

	return nil
}

func supplierRegister(args []string) (*resty.Response, error) {
	token, err := config.Token()
	if err != nil {
		logx.Logger.Error("PS: Should login first.")
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
			logx.Logger.Info("register action", "invoke expression", expVal)
			break
		}
	}

	if len(fMap) == 0 {
		logx.Logger.Info("register action", "invoke expression", "manually")
	}

	if pldVal := strings.TrimSpace(config.Vp.GetString(constant.FlagPayload.ValStr())); pldVal != "" {
		temp := make(map[string]string, 1)
		if err := json.Unmarshal([]byte(pldVal), &temp); err != nil {
			return nil, errorx.BadRequest(fmt.Sprintf("invalid payload: %s", pldVal))
		}

		fMap["payload"] = pldVal
	}
	logx.Logger.Info("register action", "invoke payload", "none")

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
