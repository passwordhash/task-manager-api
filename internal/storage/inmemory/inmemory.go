package inmemory

import (
	"context"
	"fmt"
	"sync"

	"github.com/passwordhash/task-manager-api/internal/domain"
	"github.com/passwordhash/task-manager-api/internal/storage"
	"github.com/passwordhash/task-manager-api/internal/storage/model"
)

type taskStorage struct {
	mu    sync.RWMutex
	tasks map[string]*model.Task
}

func NewTaskStorage() storage.Task {
	return &taskStorage{
		tasks: make(map[string]*model.Task),
	}
}

func (t *taskStorage) Save(_ context.Context, task domain.Task) error {
	// Maybe we should handle ctx.Done here?
	const op = "taskstorage.Save"

	t.mu.Lock()
	defer t.mu.Unlock()

	_, exists := t.tasks[task.UUID]
	if exists {
		// If the task already exists, return an error.
		// In this case it is unnecessary to check if the task is already in the storage,
		// but for contract compliance, we do it.
		return fmt.Errorf("%s: %w", op, storage.ErrAlreadyExists)
	}

	storageTask := model.FromDomainToTask(task)

	t.tasks[task.UUID] = storageTask

	return nil
}

func (t *taskStorage) Get(_ context.Context, uuid string) (domain.Task, error) {
	const op = "taskstorage.Get"

	t.mu.RLock()
	defer t.mu.RUnlock()

	task, exists := t.tasks[uuid]
	if !exists {
		return domain.Task{}, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
	}

	return task.ToDomain(), nil
}

func (t *taskStorage) Update(_ context.Context, task domain.Task) error {
	const op = "taskstorage.Update"

	t.mu.Lock()
	defer t.mu.Unlock()

	_, exists := t.tasks[task.UUID]
	if !exists {
		return fmt.Errorf("%s: %w", op, storage.ErrNotFound)
	}

	storageTask := model.FromDomainToTask(task)
	t.tasks[task.UUID] = storageTask

	return nil
}
