package task

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/passwordhash/task-manager-api/internal/domain/model"
)

type SimulatedTaskService struct {
	log *slog.Logger
}

func NewMockTaskService(log *slog.Logger) *SimulatedTaskService {
	return &SimulatedTaskService{
		log: log,
	}
}

func (m *SimulatedTaskService) CreateTask(_ context.Context) (string, error) {
	const op = "MockTaskService.CreateTask"

	log := m.log.With(slog.String("op", op))

	log.Info("Creating a mock task")

	taskUUUID := uuid.NewString()
	task := model.Task{
		UUID:      taskUUUID,
		CreatedAt: time.Now(),
		Status:    model.StatusPending,
	}

	log.Info("Mock task created", slog.String("task_uuid", task.UUID))

	// TODO: save task

	return taskUUUID, nil
}
