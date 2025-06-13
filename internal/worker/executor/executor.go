package executor

// Package executor provides an mock implementation of the Executor
// interface that simulates I/O operations.

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

func New() *simulateIOExecutor {
	return &simulateIOExecutor{}
}

func (e *simulateIOExecutor) Execute(ctx context.Context, task *domain.Task) (time.Time, error) {
	done := make(chan struct{})
	go func() {
		time.Sleep(ioDuration)
		close(done)
	}()

	select {
	case <-ctx.Done():
		return time.Now(), ctx.Err()
	case <-done:
		return time.Now(), nil
	}
}
