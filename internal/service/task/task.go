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
)

type simulatedTaskService struct {
	log         *slog.Logger
	taskStorage storage.Task
}

func NewSimulatesTaskService(
	log *slog.Logger,
	taskStorage storage.Task,
) service.TaskService {
	return &simulatedTaskService{
		log:         log,
		taskStorage: taskStorage,
	}
}

func (m *simulatedTaskService) CreateTask(ctx context.Context) (string, error) {
	const op = "MockTaskService.CreateTask"

	log := m.log.With(slog.String("op", op))

	log.Info("Creating a mock task")

	taskUUUID := uuid.NewString()
	task := domain.Task{
		UUID:      taskUUUID,
		CreatedAt: time.Now(),
		Status:    domain.StatusPending,
	}

	// TODO: move error handling to a separate function
	err := m.taskStorage.Save(ctx, task)
	if errors.Is(err, storage.ErrTaskAlreadyExist) {
		log.Error("Task with the same UUID already exists")
		return "", fmt.Errorf("%s: %w", op, service.ErrTaskAlreadyExist)
	}
	if err != nil {
		log.Error("Failed to save task", slog.Any("error", err))
		return "", fmt.Errorf("%s: failed to save task: %v", op, err)
	}

	log.Info("Mock task created and saved", "task", task)

	return taskUUUID, nil
}
