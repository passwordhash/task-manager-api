package storage

import (
	"context"
	"errors"

	"github.com/passwordhash/task-manager-api/internal/domain"
)

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound      = errors.New("not found")
)

// Task defines the interface for task storage operations.
type Task interface {
	// Save persists a task in the storage. If the task already exists,
	// it returns an [ErrAlreadyExists]. Thread safety is guaranteed.
	Save(ctx context.Context, task domain.Task) (err error)

	// Get retrieves a task by its UUID. If the task does not exist,
	// it returns an [ErrNotFound]. Thread safety is guaranteed.
	Get(ctx context.Context, uuid string) (task domain.Task, err error)

	Update(ctx context.Context, task domain.Task) (err error)
}
