package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/passwordhash/task-manager-api/internal/app"
	"github.com/passwordhash/task-manager-api/internal/config"
)

func main() {
	ctx := context.Background()

	cfg := config.MustLoad()

	log := config.SetupLogger(cfg.App.Env)

	application := app.New(ctx, log, cfg)

	go application.HTTPSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("received signal", "signal", sign)

	application.HTTPSrv.Stop()

	log.Info("stopped Task Manager API application")
}
