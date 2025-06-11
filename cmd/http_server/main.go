package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/passwordhash/task-manager-api/internal/app"
	"github.com/passwordhash/task-manager-api/internal/config"
	"golang.org/x/net/context"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "application panic: %v\n\n%s", r, debug.Stack())
			os.Exit(1)
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	cfg := config.MustLoad()

	log := config.SetupLogger(cfg.App.Env)

	application := app.New(log, cfg)

	go application.HTTPSrv.MustRun(ctx)

	<-ctx.Done()

	log.Info("received signal stop signal")

	application.HTTPSrv.Stop(ctx)

	log.Info("stopped Task Manager API application")
}
