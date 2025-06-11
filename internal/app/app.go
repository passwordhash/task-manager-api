package app

import (
	"log/slog"

	httpapp "github.com/passwordhash/task-manager-api/internal/app/http"
	"github.com/passwordhash/task-manager-api/internal/config"
	"github.com/passwordhash/task-manager-api/internal/service/task"
	"github.com/passwordhash/task-manager-api/internal/storage/inmemory"
	"github.com/passwordhash/task-manager-api/internal/worker/executor"
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

	exec := executor.New(log.WithGroup("executor"))

	workerPool := pool.New(
		log.WithGroup("pool"),
		cfg.App.Workers,
		cfg.App.TaskQueueSize,
		exec,
		taskStorage,
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
		cfg.HTTP.Port,
		cfg.HTTP.ReadTimeout,
		cfg.HTTP.WriteTimeout,
	)

	return &App{
		HTTPSrv: httpApp,
	}
}
