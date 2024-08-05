package command

import (
	"fmt"
	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/constant"
	"github.com/57blocks/auto-action/cli/internal/util"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Root represents the base command when called without any subcommands
var (
	Root = &cobra.Command{
		Use:   "aution",
		Short: "The CLI toll for auto-action: aution",
		Long: `A CLI tool: aution, which helps users to run their 
method functions on Amazon Lambda quickly, together with the result
of the execution.`,
		Args:      cobra.OnlyValidArgs,
		ValidArgs: []string{"configure", "login", "logout"},
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
			HiddenDefaultCmd:  true,
		},
		PreRun: util.PreRunBindFlags,
		Run:    rootFunc,
	}
	version = "v0.0.1"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the Root.
func Execute() {
	cobra.CheckErr(Root.Execute())
}

func init() {
	initConfig()

	Root.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
	Root.SetHelpCommand(&cobra.Command{
		Use:    "completion",
		Hidden: true,
	})
	Root.SetVersionTemplate(`Version: {{.Version}}`)
	Root.Version = version

	flagCred := constant.FlagCredential.ValStr()
	Root.PersistentFlags().StringP(
		flagCred,
		"c",
		viper.GetString(flagCred),
		"the credential file for the command about to be executed")

	flagEnv := constant.FlagEnvironment.ValStr()
	Root.PersistentFlags().StringP(
		flagEnv,
		"e",
		viper.GetString(flagEnv),
		"the execution environment")
}

func initConfig() {
	cfg, _ := config.FindOrInit()

	config.SetupViper(cfg)

	slog.Info(fmt.Sprintf("using config path: %s", viper.ConfigFileUsed()))
}

func rootFunc(cmd *cobra.Command, args []string) {
	// TODO: remove the testing code below
	fmt.Println("Root Func:")
	fmt.Println("----> viper settings:")
	for k, v := range viper.AllSettings() {
		fmt.Printf("%v: %v\n", k, v)
	}
	fmt.Println("----> args:")
	for _, v := range args {
		fmt.Printf("%v\n", v)
	}

	fmt.Println("----> flags:")
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		fmt.Printf("flag.Name: %v, flag.Value: %v\n", flag.Name, flag.Value)
	})

	cmd.Usage()
}
