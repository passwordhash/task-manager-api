package service

import (
	"context"
	"errors"

	"github.com/passwordhash/task-manager-api/internal/domain"
)

var (
	// ErrNotFound is returned when a task with the specified UUID does not exist.
	ErrNotFound = errors.New("task not found")

	// ErrAlreadyExist is returned when a task with the same UUID already exists.
	ErrAlreadyExist = errors.New("entity already exists")

	// ErrCantCancel is returned when a task cannot be canceled, for example,
	// if it is already completed or canceled.
	ErrCantCancel = errors.New("task cannot be canceled")

	ErrCantSubmit = errors.New("task cannot be submitted to worker pool")
)

// TaskService defines the interface for task-related operations.
type TaskService interface {
	// CreateTask creates a new task with status [domain.StatusPending] and returns its UUID.
	// If task cannot be submitted to the worker pool, it returns [ErrCantSubmit].
	CreateTask(ctx context.Context) (uuid string, err error)

	// Get retrieves a task by its UUID.
	// Returns [ErrNotFound] if the task does not exist.
	Get(ctx context.Context, uuid string) (task domain.Task, err error)

	// GetAll retrieves all tasks from the storage.
	GetAll(ctx context.Context) (tasks []domain.Task, err error)

	// Cancel cancels a task with the specified UUID.
	// Return [ErrNotFound] if the task does not exist,
	// [ErrCantCancel] if the task cannot be canceled
	// or some internal error.
	Cancel(ctx context.Context, uuid string) error
}
