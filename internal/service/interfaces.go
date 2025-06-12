package service

import (
	"context"
	"errors"

	"github.com/passwordhash/task-manager-api/internal/domain"
)

var (
	// ErrTaskAlreadyExist is returned when a task with the same UUID already exists.
	ErrTaskAlreadyExist = errors.New("task with the same UUID already exists")
)

// TaskService defines the interface for task-related operations.
type TaskService interface {
	// CreateTask creates a new task with status [domain.StatusPending] and returns its UUID.
	CreateTask(ctx context.Context) (uuid string, err error)

	// GetAll retrieves all tasks from the storage.
	GetAll(ctx context.Context) (tasks []domain.Task, err error)

	// TODO: doc
	Cancel(ctx context.Context, uuid string) error
}
