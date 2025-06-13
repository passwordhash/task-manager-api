package worker

import (
	"context"
	"time"

	"github.com/passwordhash/task-manager-api/internal/domain"
)

// TaskPool defines the interface for a pool of workers
// that can execute tasks concurrently.
type TaskPool interface {
	// Start initializes workers and starts processing tasks.
	Start(ctx context.Context)

	// Submit adds a task to the pool for execution.
	Submit(ctx context.Context, task *domain.Task) error

	// Cancel stops a specific task by its ID.
	Cancel(ctx context.Context, taskID string) error

	// Stop gracefully stops the pool pool, waiting for all tasks to complete
	// or the context to be done.
	Stop(ctx context.Context) error
}

type ExecuteResult struct {
	Result     any
	FinishedAt time.Time
}

// TaskExecutor defines the interface for executing tasks.
type TaskExecutor interface {
	// Execute runs i/ob-bound operation.
	// It returns the time when the task finished (even if it failed),
	// and an error if the execution failed.
	Execute(ctx context.Context, task *domain.Task) (result *ExecuteResult, error error)
}
