package task

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/passwordhash/task-manager-api/internal/domain"
	"github.com/passwordhash/task-manager-api/internal/service"
	"github.com/passwordhash/task-manager-api/internal/storage"
	"github.com/passwordhash/task-manager-api/internal/worker"
)

type simulatedTaskService struct {
	log        *slog.Logger
	workerPool worker.TaskPool
	storage    storage.Task
}

func NewSimulatedTaskService(
	log *slog.Logger,
	workerPool worker.TaskPool,
	storage storage.Task,
) service.TaskService {
	return &simulatedTaskService{
		log:        log,
		workerPool: workerPool,
		storage:    storage,
	}
}

func (m *simulatedTaskService) CreateTask(ctx context.Context) (string, error) {
	const op = "task.CreateTask"

	log := m.log.With(slog.String("op", op))

	task := domain.Task{
		UUID:      uuid.NewString(),
		CreatedAt: time.Now(),
		Status:    domain.StatusPending,
	}

	err := m.storage.Save(ctx, task)
	if err != nil {
		return "", m.handleStorageError(log, op, err)
	}

	if err := m.workerPool.Submit(ctx, &task); err != nil {
		log.Error("Failed to submit task to worker pool", slog.Any("error", err))
		return "", fmt.Errorf("%s: %w", op, service.ErrCantSubmit)
	}

	log.Info("Task created and saved", "task", task)

	return task.UUID, nil
}

func (m *simulatedTaskService) Get(ctx context.Context, uuid string) (task domain.Task, err error) {
	const op = "MockTaskService.Get"

	log := m.log.With(slog.String("op", op), slog.String("task_uuid", uuid))

	task, err = m.storage.Get(ctx, uuid)
	if err != nil {
		return domain.Task{}, m.handleStorageError(log, op, err)
	}

	log.Info("Retrieved task", slog.Any("task", task))

	return task, nil
}

func (m *simulatedTaskService) GetAll(ctx context.Context) (tasks []domain.Task, err error) {
	const op = "MockTaskService.GetAll"

	log := m.log.With(slog.String("op", op))

	tasks, err = m.storage.GetAll(ctx)
	if err != nil {
		return nil, m.handleStorageError(log, op, err)
	}

	log.Info("Retrieved all tasks", slog.Int("count", len(tasks)))

	return tasks, nil
}

func (m *simulatedTaskService) Cancel(ctx context.Context, uuid string) error {
	const op = "MockTaskService.Cancel"

	log := m.log.With(slog.String("op", op), slog.String("task_uuid", uuid))

	task, err := m.storage.Get(ctx, uuid)
	if err != nil {
		return m.handleStorageError(log, op, err)
	}

	if task.Status == domain.StatusCompleted || task.Status == domain.StatusCancelled {
		log.Warn("Task is already completed or cancelled", slog.Any("task_status", task.Status))
		return fmt.Errorf("%s: %w", op, service.ErrCantCancel)
	}

	if err := m.workerPool.Cancel(ctx, uuid); err != nil {
		log.Error("Failed to cancel task", slog.Any("error", err))
		return fmt.Errorf("%s: failed to cancel task: %v", op, err)
	}

	log.Info("Task cancelled successfully")

	return nil
}

// handleStorageError processes storage errors and returns a formatted error message.
// It checks for specific storage errors like [storage.ErrNotFound] and [storage.ErrAlreadyExists].
func (m *simulatedTaskService) handleStorageError(log *slog.Logger, op string, err error) error {
	if errors.Is(err, storage.ErrNotFound) {
		log.Warn("Task not found", slog.Any("error", err))
		return fmt.Errorf("%s: task not found: %w", op, service.ErrNotFound)
	}
	if errors.Is(err, storage.ErrAlreadyExists) {
		log.Warn("Task already exists", slog.Any("error", err))
		return fmt.Errorf("%s: task already exists: %w", op, service.ErrAlreadyExist)
	}
	log.Error("Unexpected storage error", slog.Any("error", err))
	return fmt.Errorf("%s: unexpected storage error: %v", op, err)
}
