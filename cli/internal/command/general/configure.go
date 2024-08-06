package general

import (
	"fmt"

	"github.com/57blocks/auto-action/cli/internal/command"

	"github.com/spf13/cobra"
)

// configure represents the configure command
var configure = &cobra.Command{
	Use:   "configure",
	Short: "Configure the configuration file",
	Long: `Configure the configuration file under the default path.

The config path on Mac is $HOME/.auto-action`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("configure called")
	},
}

func init() {
	command.Root.AddCommand(configure)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configure.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configure.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
