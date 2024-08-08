package oauth

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/57blocks/auto-action/cli/internal/command"
	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// logout represents the logout command
var logout = &cobra.Command{
	Use:   "logout",
	Short: "Logout the current session",
	Long: `Logout the current session under the path of credential in
the config.

For other credentials, they are still alive. It's recommended that
switching session by using **configure** command to set other 
credentials to the config`,
	Args: cobra.NoArgs,
	RunE: logoutFunc,
}

func init() {
	command.Root.AddCommand(logout)
}

func logoutFunc(_ *cobra.Command, _ []string) error {
	cfg, err := config.ReadConfig(util.DefaultPath())
	if err != nil {
		return errors.New(fmt.Sprintf("reading configuration error: %s\n", err.Error()))
	}

	if util.IsExists(cfg.Credential) {
		if err := os.Remove(cfg.Credential); err != nil {
			return errors.New(fmt.Sprintf("removing credential error: %s\n", err.Error()))
		}
	} else {
		slog.Info("credential file not found")
	}

	cfg.Credential = ""
	if err := config.WriteConfig(cfg, viper.ConfigFileUsed()); err != nil {
		return errors.New(fmt.Sprintf("writing configuration error: %s\n", err.Error()))
	}

	slog.Info("logout successfully")
	return nil
}
