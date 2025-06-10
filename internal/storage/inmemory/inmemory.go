package inmemory

import (
	"context"
	"fmt"
	"sync"

	"github.com/passwordhash/task-manager-api/internal/domain"
	"github.com/passwordhash/task-manager-api/internal/storage"
	storageModel "github.com/passwordhash/task-manager-api/internal/storage/model"
)

type taskStorage struct {
	mu    sync.RWMutex
	tasks map[string]*storageModel.Task
}

func NewTaskStorage() storage.Task {
	return &taskStorage{
		tasks: make(map[string]*storageModel.Task),
	}
}

func (t *taskStorage) Save(_ context.Context, task domain.Task) error {
	// Maybe we should handle ctx.Done here?

	const op = "storage.Save"

	t.mu.Lock()
	defer t.mu.Unlock()

	storageTask := &storageModel.Task{
		Status:    string(task.Status),
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}

	_, exists := t.tasks[task.UUID]
	if exists {
		// If the task already exists, return an error.
		// In this case it is unnecessary to check if the task is already in the storage,
		// but for contract compliance, we do it.
		return fmt.Errorf("%s: %w", op, storage.ErrTaskAlreadyExist)
	}

	t.tasks[task.UUID] = storageTask

	return nil
}
