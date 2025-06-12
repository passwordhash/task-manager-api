package executor

import (
	"context"
	"log/slog"
	"time"

	"github.com/passwordhash/task-manager-api/internal/domain"
)

const ioDuration = 100 * time.Second

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

	done := make(chan struct{})
	go func() {
		time.Sleep(ioDuration)
		close(done)
	}()

	go func() {
		for {
			select {
			case <-task.Cancel:
				log.Debug("Task execution cancelled", slog.String("task_uuid", task.UUID))
				return
			case <-done:
			}
		}
	}()

	<-done
	finishedAt := time.Now()

	log.Debug("Task execution completed", slog.String("task_uuid", task.UUID))

	return finishedAt, nil
}
