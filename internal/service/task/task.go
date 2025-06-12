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

func (m *simulatedTaskService) Cancel(ctx context.Context, uuid string) error {
	const op = "task.Cancel"

	log := m.log.With(slog.String("op", op), slog.String("uuid", uuid))

	log.Info("Cancelling task")

	// TODO: some logic ...

	if err := m.workerPool.Cancel(ctx, uuid); err != nil {
		log.Error("Failed to cancel task", slog.Any("error", err))
		return fmt.Errorf("%s: failed to cancel task: %v", op, err)
	}

	log.Info("Task cancelled successfully")

	return nil
}
