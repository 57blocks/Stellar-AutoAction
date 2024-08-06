package general

import (
	"fmt"
	"log/slog"

	"github.com/57blocks/auto-action/cli/internal/command"
	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/constant"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configure represents the configure command
var configure = &cobra.Command{
	Use:   "configure",
	Short: "Configure the configuration file",
	Long: `Configure the configuration file under the default path.

The configuration path on Mac is $HOME/.st3llar`,
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !cmd.Flags().Changed(constant.FlagCredential.ValStr()) &&
			!cmd.Flags().Changed(constant.FlagEnvPrefix.ValStr()) &&
			!cmd.Flags().Changed(constant.FlagLogLevel.ValStr()) {
			return errors.New("at least one of the flags must be set")
		}

		return nil
	},
	RunE: configureFunc,
}

func init() {
	command.Root.AddCommand(configure)

	fCred := constant.FlagCredential.ValStr()
	configure.Flags().StringP(
		fCred,
		"",
		viper.GetString(fCred),
		"configure the credential file path")

	fEnvPrefix := constant.FlagEnvPrefix.ValStr()
	configure.Flags().StringP(
		fEnvPrefix,
		"",
		viper.GetString(fEnvPrefix),
		"configure the name prefix of the environment variables")

	fLogLevel := constant.FlagLogLevel.ValStr()
	configure.Flags().StringP(
		fLogLevel,
		"",
		viper.GetString(fLogLevel),
		"configure the log level")
}

func configureFunc(_ *cobra.Command, _ []string) error {
	err := config.SyncConfig(viper.ConfigFileUsed())
	slog.Debug(fmt.Sprintf("synced configuration: %s\n", viper.ConfigFileUsed()))

	return err
}
