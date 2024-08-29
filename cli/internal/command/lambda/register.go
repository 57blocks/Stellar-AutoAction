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
	Long: `Upload/Register the local handler/handlers to Amazon Lambda, with the 
recurring/scheduled rule. 

Rules:
1. Manually, if no flags puts in, which means the handler/handlers will be triggered manually.
2. By corn
3. By rate, only three unit supported: minute, hour, day. For example: rate(1 minute).

Only cron/rate would create an event scheduler to trigger the lambda function.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed(constant.FlagCron.ValStr()) &&
			cmd.Flags().Changed(constant.FlagRate.ValStr()) {
			return errors.New("cron and rate cannot be used together")
		}

		return nil
	},
	RunE: registerFunc,
}

func init() {
	command.Root.AddCommand(register)

	flagCron := constant.FlagCron.ValStr()
	register.Flags().StringP(
		flagCron,
		"",
		viper.GetString(flagCron),
		`The cron expression for the lambda function. Which is used in the Event Scheduler.
For more info: https://docs.aws.amazon.com/AmazonCloudWatch/latest/events/ScheduledEvents.html`)

	flagRate := constant.FlagRate.ValStr()
	register.Flags().StringP(
		flagRate,
		"",
		viper.GetString(flagRate),
		`The cron expression for the lambda function. Which is used in the Event Scheduler.
For more info: https://docs.aws.amazon.com/AmazonCloudWatch/latest/events/ScheduledEvents.html`)
}

func registerFunc(cmd *cobra.Command, args []string) error {
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

	if len(args) > 0 {
		flagMap := make(map[string]string, 1)
		if strings.TrimSpace(viper.GetString(constant.FlagCron.ValStr())) == "" {
			slog.Debug(fmt.Sprintf("rate expression: %s\n", viper.GetString(constant.FlagRate.ValStr())))
			flagMap["expression"] = viper.GetString(constant.FlagRate.ValStr())
		} else {
			slog.Debug(fmt.Sprintf("cron expression: %s\n", viper.GetString(constant.FlagCron.ValStr())))
			flagMap["expression"] = viper.GetString(constant.FlagCron.ValStr())
		}

		request = request.SetFormData(flagMap)
	}

	resp, err := request.Post(viper.GetString("bound_with.endpoint") + "/lambda/register")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("resty error: %s\n", err.Error()))
	}

	if e := util.HasError(resp); e != nil {
		return nil, errors.New(fmt.Sprintf("supplier error: %s\n", e))
	}

	return resp, nil
}
