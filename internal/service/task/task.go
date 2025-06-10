package task

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/passwordhash/task-manager-api/internal/domain"
	"github.com/passwordhash/task-manager-api/internal/service"
)

type simulatedTaskService struct {
	log *slog.Logger
}

func NewMockTaskService(log *slog.Logger) service.TaskService {
	return &simulatedTaskService{
		log: log,
	}
}

func (m *simulatedTaskService) CreateTask(_ context.Context) (string, error) {
	const op = "MockTaskService.CreateTask"

	log := m.log.With(slog.String("op", op))

	log.Info("Creating a mock task")

	taskUUUID := uuid.NewString()
	task := domain.Task{
		UUID:      taskUUUID,
		CreatedAt: time.Now(),
		Status:    domain.StatusPending,
	}

	log.Info("Mock task created", slog.String("task_uuid", task.UUID))

	// TODO: save task

	return taskUUUID, nil
}
