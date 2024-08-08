package oauth

import (
	"fmt"
	"os"
	"syscall"

	"github.com/57blocks/auto-action/cli/internal/command"
	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/constant"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
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
	RunE:      loginFunc,
}

func init() {
	command.Root.AddCommand(login)

	flagCred := constant.FlagCredential.ValStr()
	login.Flags().StringP(
		flagCred,
		"c",
		viper.GetString(flagCred),
		"the credential file for the command about to be executed")

	flagEnv := constant.FlagEnvironment.ValStr()
	login.Flags().StringP(
		flagEnv,
		"e",
		viper.GetString(flagEnv),
		"the execution environment")

	flagAcc := constant.FlagAccount.ValStr()
	login.Flags().StringP(
		flagAcc,
		"a",
		"",
		"name of the account")

	flagOrg := constant.FlagOrganization.ValStr()
	login.Flags().StringP(flagOrg,
		"o",
		"",
		"name of the organization")

	if err := login.MarkFlagRequired(flagAcc); err != nil {
		return
	}
	if err := login.MarkFlagRequired(flagOrg); err != nil {
		return
	}
}

func loginFunc(cmd *cobra.Command, args []string) error {
	fmt.Println("Password: ")

	passwordBytes, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return errors.New(fmt.Sprintf("reading passowrd error: %s\n", err.Error()))
	}

	if len(passwordBytes) == 0 {
		return errors.New("empty password error")
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

	return syncCred()
}

func syncCred() error {
	// TODO: implement the generation logic
	session := config.BuildCredSession(
		config.WithSessionToken("JWT_TOKEN"),
		config.WithSessionRefreshToken("JWT_REFRESH"),
	)

	// TODO: put it to env
	var endpoint = "http://st3llar-alb-365211.us-east-2.elb.amazonaws.com"
	env := config.BuildCredEnv(
		config.WithEnvName(viper.GetString(constant.FlagEnvironment.ValStr())),
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
		return errors.New(fmt.Sprintf("marshaling credentials error: %s\n", err.Error()))
	}

	// Write TOML to credential file
	err = os.WriteFile(viper.GetString(constant.FlagCredential.ValStr()), credToml, 0666)
	if err != nil {
		return errors.New(fmt.Sprintf("writing credentials error: %s\n", err.Error()))
	}

	// sync to configuration file
	return config.SyncConfig(viper.ConfigFileUsed())
}
