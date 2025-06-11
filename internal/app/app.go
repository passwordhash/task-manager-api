package app

import (
	"log/slog"

	httpapp "github.com/passwordhash/task-manager-api/internal/app/http"
	"github.com/passwordhash/task-manager-api/internal/config"
	"github.com/passwordhash/task-manager-api/internal/service/task"
	"github.com/passwordhash/task-manager-api/internal/storage/inmemory"
	"github.com/passwordhash/task-manager-api/internal/worker/pool"
)

type App struct {
	HTTPSrv *httpapp.App
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	taskStorage := inmemory.NewTaskStorage()

	workerPool := pool.New(
		log.WithGroup("pool"),
		100, // TODO: make configurable
		3,   // TODO: make configurable
		nil, // TODO: make configurable
	)

	taskService := task.NewSimulatedTaskService(
		log.WithGroup("service"),
		workerPool,
		taskStorage,
	)

	httpApp := httpapp.New(
		log,
		workerPool,
		taskService,
		cfg.App.Port,
		cfg.App.ReadTimeout,
		cfg.App.WriteTimeout,
	)

	return &App{
		HTTPSrv: httpApp,
	}
}
