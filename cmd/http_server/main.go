package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/passwordhash/task-manager-api/internal/app"
	"github.com/passwordhash/task-manager-api/internal/config"
	"golang.org/x/net/context"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	cfg := config.MustLoad()

	log := config.SetupLogger(cfg.Env)

	application := app.New(ctx, log, cfg)

	go application.HTTPSrv.MustRun(ctx)

	<-ctx.Done()

	log.Info("received signal stop signal")

	application.HTTPSrv.Stop(ctx)

	log.Info("stopped Task Manager API application")
}
