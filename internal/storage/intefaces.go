package storage

import (
	"context"
	"errors"

	"github.com/passwordhash/task-manager-api/internal/domain"
)

var (
	ErrTaskAlreadyExist = errors.New("task already exists")
)

// Task defines the interface for task storage operations.
type Task interface {
	// Save persists a task in the storage. If the task already exists,
	// it returns an [ErrTaskAlreadyExist].
	Save(ctx context.Context, task domain.Task) (err error)
}
