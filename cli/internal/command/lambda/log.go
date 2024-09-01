package lambda

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/57blocks/auto-action/cli/internal/command"
	"github.com/57blocks/auto-action/cli/internal/config"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// logs represents the log command
var logs = &cobra.Command{
	Use:   "log",
	Short: "Track logs of the lambda function",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return logFunc(cmd, args)
	},
}

func init() {
	command.Root.AddCommand(logs)
}

func logFunc(cmd *cobra.Command, args []string) error {
	fmt.Println(viper.GetString("bound_with.endpoint"))

	token, err := config.Token()
	if err != nil {
		return err
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/lambda/logs/" + args[0]}
	slog.Info(fmt.Sprintf("connecting to %s", u.String()))

	// Add JWT token to the request headers
	header := http.Header{}
	header.Add("Authorization", token)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
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
				log.Println("write:", err)
				return err
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return err
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}

//response, err := restyx.Client.R().
//	EnableTrace().
//	SetHeaders(map[string]string{
//		"Content-Type":  "multipart/form-data",
//		"Authorization": token,
//	}).
//	Get(viper.GetString("bound_with.endpoint") + "/lambda/logs/" + args[0])
//if err != nil {
//	return errors.Wrap(err, fmt.Sprintf("resty error: %s\n", err.Error()))
//}
//
//slog.Debug(fmt.Sprintf("%v\n", response)) // TODO: remove
//
//if e := util.HasError(response); e != nil {
//	return errors.Wrap(e, fmt.Sprintf("supplier error: %s\n", e))
//}

//	return nil
//}
