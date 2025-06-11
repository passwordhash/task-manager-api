package worker

import (
	"context"

	"github.com/passwordhash/task-manager-api/internal/domain"
)

// TaskPool defines the interface for a pool of workers
// that can execute tasks concurrently.
type TaskPool interface {
	// Start initializes workers and starts processing tasks.
	Start(ctx context.Context)

	// Submit adds a task to the pool for execution.
	Submit(ctx context.Context, task *domain.Task) error

	// Stop gracefully stops the worker pool, waiting for all tasks to complete
	// or the context to be done.
	Stop(ctx context.Context) error
}
