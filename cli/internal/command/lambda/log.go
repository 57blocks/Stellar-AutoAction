package lambda

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/57blocks/auto-action/cli/internal/command"
	"github.com/57blocks/auto-action/cli/internal/config"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// logs represents the log command
var logs = &cobra.Command{
	Use:   "log <name/arn>",
	Short: "Tracking execution logs of the lambda function",
	Long: `
Tracking execution logs of a specific Lambda, by its name/arn.

Which with a 5 seconds interval to fetch the logs.
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return logFunc(cmd, args)
	},
}

func init() {
	command.Root.AddCommand(logs)
}

func logFunc(_ *cobra.Command, args []string) error {
	token, err := config.Token()
	if err != nil {
		return err
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	splits := strings.Split(viper.GetString("bound_with.endpoint"), "://")

	u := url.URL{
		Scheme: "ws",
		Host:   splits[1],
		Path:   fmt.Sprintf("/lambda/%s/logs", args[0]),
	}
	slog.Info(fmt.Sprintf("dailing to %s\n", u.String()))

	// Add JWT token to the request headers
	header := http.Header{}
	header.Add("Authorization", token)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		slog.Error(fmt.Sprintf("dialing error: %s\n", err.Error()))
		return errors.Wrap(err, "failed to dial to websocket")
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				slog.Error(fmt.Sprintf("read cloudwatch events error: %s", err.Error()))

				return
			}

			slog.Info(fmt.Sprintf("%s", message))
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return nil
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				return errors.Wrap(err, "failed to write message to websocket")
			}
		case <-interrupt:
			slog.Info("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return errors.Wrap(err, "failed to write close message to websocket")
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}
