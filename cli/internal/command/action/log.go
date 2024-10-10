package action

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/pkg/errorx"
	"github.com/57blocks/auto-action/cli/internal/pkg/logx"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

// logs represents the `log` command
var logs = &cobra.Command{
	Use:   "log <name/arn>",
	Short: "Tracking execution logs of the action",
	Long: `
Description:
  Tracking execution logs of a specific Action, by name/arn.
  Which with a 5 seconds interval to fetch the latest events.

TODO:
  - Add time range filer
  - Add error events filer
`,
	Args: cobra.ExactArgs(1),
	RunE: logFunc,
}

func init() {
	actionGroup.AddCommand(logs)
}

func logFunc(_ *cobra.Command, args []string) error {
	token, err := config.Token()
	if err != nil {
		return err
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	splits := strings.Split(config.Vp.GetString("bound_with.endpoint"), "://")

	u := url.URL{
		Scheme: "ws",
		Host:   splits[1],
		Path:   util.ParseReqPath(fmt.Sprintf("/lambda/%s/logs", args[0])),
	}
	logx.Logger.Debug("ws", "dailing to", u.String())

	// Add JWT token to the request headers
	header := http.Header{}
	header.Add("Authorization", token)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		logx.Logger.Error("ws", "dialing error", err.Error())
		return errorx.Internal(err.Error())
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				logx.Logger.Error("ws", "read cloudwatch events error", err.Error())

				return
			}

			logx.Logger.Info("CloudWatch Event", "detail", string(message))
		}
	}()

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return nil
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				return errorx.Internal(fmt.Sprintf("failed to write message to websocket: %s", err.Error()))
			}
		case <-interrupt:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return errorx.Internal(fmt.Sprintf("failed to write close message to websocket: %s", err.Error()))
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}
