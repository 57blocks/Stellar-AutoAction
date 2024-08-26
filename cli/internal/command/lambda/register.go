package lambda

import (
	"fmt"
	"log/slog"

	"github.com/57blocks/auto-action/cli/internal/command"
	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/pkg/restyx"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// register represents the register command
var register = &cobra.Command{
	Use:   "register",
	Short: "Lambda register of local handler",
	Long: `Upload/Register the local handler/handlers to Amazon Lambda, with the 
recurring/scheduled rule. 

Rules:
1. 

And if not, the handler/handlers will be triggered manually in further future.
Finally, returns the lambda ARN`,
	RunE: registerFunc,
}

func init() {
	command.Root.AddCommand(register)
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
	argMap := make(map[string]string, len(args))
	for _, path := range args {
		argMap[path] = path
	}

	token, err := config.Token()
	if err != nil {
		return nil, err
	}

	resp, err := restyx.Client.R().
		EnableTrace().
		SetHeaders(map[string]string{
			"Content-Type":  "multipart/form-data",
			"Authorization": token,
		}).
		SetFiles(argMap).
		Post(viper.GetString("bound_with.endpoint") + "/lambda/register")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("resty error: %s\n", err.Error()))
	}

	if e := util.HasError(resp); e != nil {
		return nil, errors.New(fmt.Sprintf("supplier error: %s\n", e))
	}

	return resp, nil
}
