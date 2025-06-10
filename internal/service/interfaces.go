package service

import (
	"context"

	"github.com/passwordhash/task-manager-api/internal/domain"
)

type TaskService interface {
	CreateTask(ctx context.Context, task *domain.Task) (uuid string, err error)
}
