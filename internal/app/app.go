package app

import (
	"context"
	"log/slog"

	httpapp "github.com/passwordhash/task-manager-api/internal/app/http"
	"github.com/passwordhash/task-manager-api/internal/config"
	"github.com/passwordhash/task-manager-api/internal/services/task"
)

type App struct {
	HTTPSrv *httpapp.App
}

func New(
    ctx context.Context,
    log *slog.Logger,
    cfg *config.Config,
) *App {
	taskService := task.NewMockTaskService(log)

	httpApp := httpapp.New(
		ctx,
		log,
		taskService,
		cfg.App.Port,
		cfg.App.WriteTimeout,
		cfg.App.ReadTimeout,
	)

	return &App{
		HTTPSrv: httpApp,
	}
}
