package executor

import (
	"context"
	"log/slog"
	"time"

	"github.com/passwordhash/task-manager-api/internal/domain"
)

const ioDuration = 5 * time.Second

type simulateIOExecutor struct {
	log *slog.Logger
}

func New(log *slog.Logger) *simulateIOExecutor {
	return &simulateIOExecutor{
		log: log,
	}
}

func (e *simulateIOExecutor) Execute(_ context.Context, task *domain.Task) (time.Time, error) {
	const op = "executor.Execute"

	log := e.log.With(slog.String("op", op), slog.String("task_uuid", task.UUID))

	log.Debug("Starting task execution")

	time.Sleep(ioDuration)
	finishedAt := time.Now()

	log.Debug("Task execution completed", slog.String("task_uuid", task.UUID))

	return finishedAt, nil
}
