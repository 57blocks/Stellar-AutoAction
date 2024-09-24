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
	Short: "Configure the configuration file",
	Long: `
Description:
  Configure the configuration file under the default path.

Note:
  - When specifying other credentials, please confirm with that the
    credential is matched with the bound endpoint and not expired.
  - When specifying the log level, here are the options below:
    - Debug
    - Warn
    - Error
    - Info
    If none matched, using **Info** as default.
  - When specifying the tracking source, here are the options below:
    - ON
    - OFF
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
		"configure the credential file path")

	fEndPoint := constant.FlagEndPoint.ValStr()
	configure.Flags().StringP(
		fEndPoint,
		"",
		config.Vp.GetString(fEndPoint),
		"configure the endpoint of the service")

	fLogLevel := constant.FlagLog.ValStr()
	configure.Flags().StringP(
		fLogLevel,
		"",
		config.Vp.GetString(fLogLevel),
		"configure the logx level")

	fSource := constant.FlagSource.ValStr()
	configure.Flags().StringP(
		fSource,
		"s",
		config.Vp.GetString(fSource),
		"configure the tracking source or not")
}

func configureFunc(_ *cobra.Command, _ []string) error {
	if err := config.SyncConfigByFlags(); err != nil {
		return err
	}

	logx.Logger.Info(fmt.Sprintf("synced config: %s", config.Vp.ConfigFileUsed()))

	return nil
}
