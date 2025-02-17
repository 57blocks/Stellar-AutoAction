package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/57blocks/auto-action/server/internal/api"
	"github.com/57blocks/auto-action/server/internal/boot"
	"github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/db"
	"github.com/57blocks/auto-action/server/internal/service"
	thirdParty "github.com/57blocks/auto-action/server/internal/third-party"
	"github.com/57blocks/auto-action/server/internal/third-party/logx"
)

var server *http.Server

func main() {
	if err := boot.Boots(
		boot.Wrap(func() error {
			return config.Setup("./internal/config/")
		}),
		boot.Wrap(func() error {
			return logx.Setup(config.GlobalConfig)
		}),
		boot.Wrap(db.Setup),
		boot.Wrap(thirdParty.Setup),
		boot.Wrap(service.Setup),
		boot.Wrap(api.Setup),
	); err != nil {
		log.Panicf("boots components occurred error: %s", err.Error())
	}

	logx.Logger.INFO("boots: server")

	server = &http.Server{
		Addr:    ":8080",
		Handler: api.GinEngine,
	}

	go server.ListenAndServe()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	if err := stopServer(); err != nil {
		log.Fatal("shutting down occurred error: ", err)
	}
	log.Println("Exited")
}

func stopServer() error {
	shutErr := make(chan error)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		shutErr <- server.Shutdown(ctx)
	}()

	select {
	case err := <-shutErr:
		if err != nil {
			log.Fatal("Server Shutdown:", err)
			return err
		}
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
		return ctx.Err()
	}

	return nil
}
