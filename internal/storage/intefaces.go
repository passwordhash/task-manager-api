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

type TaskUpdate struct {
	Status    domain.TaskStatus
	UpdatedAt time.Time
	StartedAt time.Time
	Result    any
	Error     error
}

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

	Update(ctx context.Context, uuid string, update TaskUpdate) (err error)
}
