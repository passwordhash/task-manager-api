package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	httpapp "github.com/passwordhash/task-manager-api/internal/app/http"
	"github.com/passwordhash/task-manager-api/internal/config"
	"github.com/passwordhash/task-manager-api/internal/domain"
	"github.com/passwordhash/task-manager-api/internal/service/task"
	"github.com/passwordhash/task-manager-api/internal/storage/inmemory"
	"github.com/passwordhash/task-manager-api/internal/worker"
)

type App struct {
	HTTPSrv *httpapp.App
}

type tmpExecut struct{}

func (t *tmpExecut) Execute(ctx context.Context, task *domain.Task) error {
	fmt.Println("start execution...")
	time.Sleep(5 * time.Second)
	fmt.Println("task executed successfully")

	return nil
}

func New(
	ctx context.Context,
	log *slog.Logger,
	cfg *config.Config,
) *App {
	taskStorage := inmemory.NewTaskStorage()

	workerPool := worker.NewPool(
		log.WithGroup("worker"),
		100,          // TODO: make configurable
		3,            // TODO: make configurable
		&tmpExecut{}, // TODO: make configurable
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
