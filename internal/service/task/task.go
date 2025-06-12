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

	log.Info("Creating a mock task")

	taskUUUID := uuid.NewString()
	task := domain.Task{
		UUID:      taskUUUID,
		CreatedAt: time.Now(),
		Status:    domain.StatusPending,
	}

	// TODO: move error handling to a separate function
	err := m.storage.Save(ctx, task)
	if errors.Is(err, storage.ErrAlreadyExists) {
		log.Error("Task with the same UUID already exists")
		return "", fmt.Errorf("%s: %w", op, service.ErrTaskAlreadyExist)
	}
	if err != nil {
		log.Error("Failed to save task", slog.Any("error", err))
		return "", fmt.Errorf("%s: failed to save task: %v", op, err)
	}

	m.workerPool.Submit(ctx, &task)

	log.Info("Mock task created and saved", "task", task)

	return taskUUUID, nil
}

func (m *simulatedTaskService) GetAll(ctx context.Context) (tasks []domain.Task, err error) {
	const op = "MockTaskService.GetAll"

	log := m.log.With(slog.String("op", op))

	log.Info("Retrieving all mock tasks")

	tasks, err = m.storage.GetAll(ctx)
	if err != nil {
		log.Error("Failed to retrieve tasks", slog.Any("error", err))
		return nil, fmt.Errorf("%s: failed to retrieve tasks: %v", op, err)
	}

	log.Info("Retrieved all mock tasks successfully", slog.Int("count", len(tasks)))

	return tasks, nil
}

func (m *simulatedTaskService) Cancel(ctx context.Context, uuid string) error {
	const op = "MockTaskService.Cancel"

	log := m.log.With(slog.String("op", op), slog.String("task_uuid", uuid))

	log.Info("Cancelling mock task")

	task, err := m.storage.Get(ctx, uuid)
	if err != nil {
		log.Error("Failed to retrieve task", slog.Any("error", err))
		return fmt.Errorf("%s: failed to retrieve task: %v", op, err)
	}

	if task.Status == domain.StatusCompleted || task.Status == domain.StatusCancelled {
		log.Warn("Task is already completed or cancelled", slog.Any("task_status", task.Status))
		return nil
	}

	task.Cancel <- struct{}{}

	log.Info("Mock task cancelled successfully")

	return nil
}
