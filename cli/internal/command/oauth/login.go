package oauth

import (
	"fmt"
	"os"
	"syscall"

	"github.com/57blocks/auto-action/cli/internal/command"
	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/constant"
	"github.com/57blocks/auto-action/cli/internal/util"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

// login represents the login command
var login = &cobra.Command{
	Use:   "login",
	Short: "Login to the Stellar auto-action.",
	Long: `Login the Stellar auto-action based on credential path in
the config. And will create a new credential under the path you just claimed
and set it to config, if it's the first time.

And also, you could specify other credentials by **configure** command.`,
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{"organization", "account", "environment", "credential"},
	PreRun: func(cmd *cobra.Command, args []string) {
		util.PreRunBindFlags(cmd, args)
	},
	Run: loginFunc,
}

func init() {
	command.Root.AddCommand(login)

	login.Flags().StringP("account", "A", "", "name of the account")
	login.Flags().StringP("organization", "O", "", "name of the organization")

	if err := login.MarkFlagRequired("account"); err != nil {
		return
	}
	if err := login.MarkFlagRequired("organization"); err != nil {
		return
	}
}

func loginFunc(cmd *cobra.Command, args []string) {
	fmt.Print("Password: ")
	passwordBytes, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reading passowrd occurred an error: %s\n", err.Error())
		os.Exit(1)
		return
	}

	password := string(passwordBytes)
	fmt.Printf("input password: %v\n", password)

	fmt.Println("Login Func:")
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

	cobra.CheckErr(syncCred)
}

func syncCred() error {
	session := config.BuildCredSession(
		config.WithSessionToken("JWT_TOKEN"),
		config.WithSessionRefreshToken("JWT_REFRESH"),
	)

	// TODO: put it to env
	var endpoint = "http://st3llar-alb-365211.us-east-2.elb.amazonaws.com"
	env := config.BuildCredEnv(
		config.WithEnvName("Horizon"),
		config.WithEnvEndPoint(endpoint),
	)

	cred := config.BuildCred(
		config.WithAccount(viper.GetString(constant.FlagAccount.ValStr())),
		config.WithOrganization(viper.GetString(constant.FlagOrganization.ValStr())),
		config.WithEnvironment(env),
		config.WithSession(session),
	)

	credToml, err := toml.Marshal(cred)
	if err != nil {
		_, e := fmt.Fprintf(
			os.Stderr,
			"marshaling credentials to TOML error: %s\n",
			err.Error(),
		)
		return e
	}

	// Write TOML to credential file
	err = os.WriteFile(viper.GetString(constant.FlagCredential.ValStr()), credToml, 0666)
	if err != nil {
		_, e := fmt.Fprintf(
			os.Stderr,
			"writing credentials to file error: %s\n",
			err.Error(),
		)
		return e
	}

	// sync to configuration file
	cfg := config.Build(
		config.WithCredential(viper.GetString(constant.FlagCredential.ValStr())),
		config.WithEnvVarPrefix(viper.GetString(constant.FlagEnvPrefix.ValStr())),
		config.WithLogLevel(viper.GetString(constant.FlagLogLevel.ValStr())),
	)

	if err := config.WriteConfig(cfg, util.DefaultPath()); err != nil {
		_, e := fmt.Fprintf(
			os.Stderr,
			"writing configuration to TOML error: %s\n",
			err.Error(),
		)
		return e
	}

	return nil
}
