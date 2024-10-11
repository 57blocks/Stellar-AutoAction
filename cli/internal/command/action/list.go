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
	Short: "List all actions registered in Stellar AutoAction",
	Long: `
Description:
  The list command displays all actions currently registered in Stellar AutoAction.
  It provides an overview of each action's basic details.

Examples:
  autoaction action list
  autoaction action list -f

Output:
  By default, the command output includes the function name, ARN, and creation date.

When using the --full flag, additional details are displayed:
  - Action configuration (e.g., handler, runtime, role)
  - Code information (SHA256, revision)
  - Scheduler details (if applicable)
  - Timestamps and version

Note:
  Use the -f or --full flag to retrieve more comprehensive information about
  each action, including detailed configurations and associated schedulers.
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
		`Display comprehensive details for all actions.
When enabled, this flag provides:
- Detailed action configurations
- Associated scheduler information
- Code details and timestamps
Use this for a more in-depth view of your actions.`)
}

func listFunc(cmd *cobra.Command, _ []string) error {
	token, err := config.Token()
	if err != nil {
		logx.Logger.Error("PS: Should login first.")
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
