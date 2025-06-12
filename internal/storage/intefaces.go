package storage

import (
	"context"
	"errors"
	"time"

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

	// GetAll retrieves all tasks from the storage. Thread safety is guaranteed.
	GetAll(ctx context.Context) (tasks []domain.Task, err error)

	UpdateStatus(ctx context.Context,
	    uuid string,
	    status domain.TaskStatus,
	    updatedAt time.Time,
	) error
}
