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

	// ErrCantBeCancelled is returned when a task cannot be cancelled, for example,
	// if it is already completed or cancelled.
	ErrCantBeCancelled = errors.New("task cannot be cancelled")
)

// TaskService defines the interface for task-related operations.
type TaskService interface {
	// CreateTask creates a new task with status [domain.StatusPending] and returns its UUID.
	CreateTask(ctx context.Context) (uuid string, err error)

	// GetAll retrieves all tasks from the storage.
	GetAll(ctx context.Context) (tasks []domain.Task, err error)

	// Cancel cancels a task with the specified UUID.
	// Return ErrNotFound if the task does not exist,
	// ErrCantBeCancelled if the task cannot be cancelled
	// or some internal error.
	Cancel(ctx context.Context, uuid string) error
}
