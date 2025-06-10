package app

import (
	"context"
	"log/slog"

	httpapp "github.com/passwordhash/task-manager-api/internal/app/http"
	"github.com/passwordhash/task-manager-api/internal/config"
	"github.com/passwordhash/task-manager-api/internal/service/task"
	"github.com/passwordhash/task-manager-api/internal/storage/inmemory"
)

type App struct {
	HTTPSrv *httpapp.App
}

func New(
	ctx context.Context,
	log *slog.Logger,
	cfg *config.Config,
) *App {
	taskStorage := inmemory.NewTaskStorage()

	taskService := task.NewSimulatesTaskService(log, taskStorage)

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
