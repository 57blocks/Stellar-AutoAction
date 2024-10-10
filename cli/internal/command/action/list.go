package action

import (
	"encoding/json"
	"fmt"

	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/constant"
	"github.com/57blocks/auto-action/cli/internal/pkg/errorx"
	"github.com/57blocks/auto-action/cli/internal/pkg/logx"
	"github.com/57blocks/auto-action/cli/internal/pkg/restyx"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all actions registered in AutoAction",
	Long: `
Description:
  List all actions registered in AutoAction with their details as much 
  as possible.

Note:
  Using -f/--full flag to list with more details of Action itself,
  together with their schedulers info.
`,
	Args: cobra.NoArgs,
	RunE: listFunc,
}

func init() {
	actionGroup.AddCommand(listCmd)

	fFull := constant.FlagFull.ValStr()
	listCmd.Flags().BoolP(
		fFull,
		"f",
		config.Vp.GetBool(fFull),
		`
List all actions with their details as much as possible.
`)
}

func listFunc(cmd *cobra.Command, _ []string) error {
	token, err := config.Token()
	if err != nil {
		return err
	}

	rawURL := fmt.Sprintf("%s/lambda", config.Vp.GetString("bound_with.endpoint"))
	if cmd.Flags().Lookup("full").Changed {
		rawURL = fmt.Sprintf("%s?full=true", rawURL)
	}

	URL := util.ParseReqPath(rawURL)

	response, err := restyx.Client.R().
		EnableTrace().
		SetHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": token,
		}).
		Get(URL)
	if err != nil {
		return errorx.RestyError(err.Error())
	}
	if response.IsError() {
		return errorx.WithRestyResp(response)
	}

	var respData []map[string]interface{}
	if err := json.Unmarshal(response.Body(), &respData); err != nil {
		logx.Logger.Error("Error unmarshalling JSON", "error", err.Error())
		return errorx.Internal(err.Error())
	}

	logx.Logger.Info("actions", "results", respData)

	return nil
}
