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
	"github.com/57blocks/auto-action/server/internal/pkg/db"
	pkgLog "github.com/57blocks/auto-action/server/internal/pkg/log"
)

var server *http.Server

func main() {
	if err := boot.Boots(
		boot.Wrap(config.Setup),
		boot.Wrap(pkgLog.Setup),
		boot.Wrap(db.Setup),
		boot.Wrap(api.Setup),
	); err != nil {
		log.Panicf("boots components occurred error: %s\n", err.Error())
	}

	pkgLog.Logger.INFO("boots: server")

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
