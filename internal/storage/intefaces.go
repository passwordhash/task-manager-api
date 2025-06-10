package storage

import (
	"context"

	"github.com/passwordhash/task-manager-api/internal/domain"
)

// Task defines the interface for task storage operations.
type Task interface {
	// Save persists a task in the storage. If the task already exists,
	// it returns an [ErrTaskAlreadyExist].
	Save(ctx context.Context, task domain.Task) (err error)
}
