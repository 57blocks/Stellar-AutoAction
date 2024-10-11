package general

import (
	"fmt"

	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/constant"
	"github.com/57blocks/auto-action/cli/internal/pkg/errorx"
	"github.com/57blocks/auto-action/cli/internal/pkg/logx"

	"github.com/spf13/cobra"
)

// configure represents the configure command
var configure = &cobra.Command{
	Use:   "configure",
	Short: "Modify the Stellar AutoAction configuration settings",
	Long: `
Description:
  The configure command allows you to modify various settings in the Stellar AutoAction 
  configuration file. This file is located at the default path and contains important 
  parameters that control the behavior of the Stellar AutoAction tool.

Configurable Settings:
  - Credentials: Specify alternative authentication credentials
  - Log Level: Set the verbosity of logging output
  - Tracking Source: Enable or disable action tracking

Notes:
  1. Credentials:
     - Ensure that any specified credential matches the bound endpoint and is not expired.
     - Use the full path to the credential file when specifying an alternative.

  2. Log Level:
     - Available options: Debug, Warn, Error, Info
     - If an invalid option is provided, the default level "Info" will be used.

  3. Tracking Source:
     - Available options: ON, OFF
     - If an invalid option is provided, the default "OFF" will be used.

Examples:
  autoaction configure --credential /path/to/cred.json --log Debug --source ON
`,
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !cmd.Flags().Changed(constant.FlagCredential.ValStr()) &&
			!cmd.Flags().Changed(constant.FlagEndPoint.ValStr()) &&
			!cmd.Flags().Changed(constant.FlagSource.ValStr()) &&
			!cmd.Flags().Changed(constant.FlagLog.ValStr()) {
			return errorx.BadRequest("at least one of the flags must be set")
		}

		return nil
	},
	RunE: configureFunc,
}

func init() {
	generalGroup.AddCommand(configure)

	fCred := constant.FlagCredential.ValStr()
	configure.Flags().StringP(
		fCred,
		"",
		config.Vp.GetString(fCred),
		`Path to the credential file for authentication.
Use this to specify an alternative credential file.`)

	fEndPoint := constant.FlagEndPoint.ValStr()
	configure.Flags().StringP(
		fEndPoint,
		"",
		config.Vp.GetString(fEndPoint),
		`URL of the Stellar AutoAction service endpoint.
Specify the server address to connect to.`)

	fLogLevel := constant.FlagLog.ValStr()
	configure.Flags().StringP(
		fLogLevel,
		"",
		config.Vp.GetString(fLogLevel),
		`Set the logging level (Debug, Warn, Error, Info).
Controls the verbosity of log output. Default is Info.`)

	fSource := constant.FlagSource.ValStr()
	configure.Flags().StringP(
		fSource,
		"s",
		config.Vp.GetString(fSource),
		`Enable or disable action tracking (ON, OFF).
Determines if actions are tracked. Default is OFF.`)
}

func configureFunc(_ *cobra.Command, _ []string) error {
	if err := config.SyncConfigByFlags(); err != nil {
		return err
	}

	logx.Logger.Info(fmt.Sprintf("synced config: %s", config.Vp.ConfigFileUsed()))

	return nil
}
