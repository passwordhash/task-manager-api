package executor

import (
	"context"
	"flag"
	"log/slog"
	"time"

	"github.com/passwordhash/task-manager-api/internal/domain"
)

var ioDuration time.Duration

func init() {
	flag.DurationVar(&ioDuration, "io-duration", 10*time.Second, "Duration to simulate I/O operation in the executor")
}

type simulateIOExecutor struct {
	log *slog.Logger
}

func New(log *slog.Logger) *simulateIOExecutor {
	return &simulateIOExecutor{
		log: log,
	}
}

func (e *simulateIOExecutor) Execute(ctx context.Context, task *domain.Task) (time.Time, error) {
	const op = "executor.Execute"

	log := e.log.With(slog.String("op", op), slog.String("task_uuid", task.UUID))

	log.Debug("Starting task execution")

	done := make(chan struct{})
	go func() {
		time.Sleep(ioDuration)
		close(done)
	}()

	select {
	case <-ctx.Done():
		log.Debug("Task execution cancelled", slog.String("task_uuid", task.UUID))
		return time.Time{}, ctx.Err()
	case <-done:
		finishedAt := time.Now()

		log.Debug("Task execution completed", slog.String("task_uuid", task.UUID))

		return finishedAt, nil
	}
}
