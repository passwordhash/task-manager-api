package executor

// Package executor provides an mock implementation of the Executor
// interface that simulates I/O operations.

import (
	"context"
	"flag"
	"log/slog"
	"time"

	"github.com/passwordhash/task-manager-api/internal/domain"
	"github.com/passwordhash/task-manager-api/internal/worker"
)

var ioDuration time.Duration

func init() {
	flag.DurationVar(&ioDuration, "io-duration", 10*time.Second, "Duration to simulate I/O operation in the executor")
}

type simulateIOExecutor struct {
	log *slog.Logger
}

func New() *simulateIOExecutor {
	return &simulateIOExecutor{}
}

func (e *simulateIOExecutor) Execute(ctx context.Context, task *domain.Task) (*worker.ExecuteResult, error) {
	var execRes worker.ExecuteResult

	done := make(chan struct{})
	go func() {
		time.Sleep(ioDuration)
		result := map[string]any{
			"message":  "I/O operation completed",
			"bytes":    1024,
			"duration": ioDuration.String(),
			"task":     task.UUID,
		}
		execRes.FinishedAt = time.Now()
		execRes.Result = result
		close(done)
	}()

	select {
	case <-ctx.Done():
		return &execRes, ctx.Err()
	case <-done:
		return &execRes, nil
	}
}
