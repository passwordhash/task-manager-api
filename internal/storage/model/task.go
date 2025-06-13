package model

import (
	"time"

	"github.com/passwordhash/task-manager-api/internal/domain"
)

type Task struct {
	Status    string
	CreatedAt time.Time
	StartedAt time.Time
	UpdatedAt time.Time
	Result    any
	Error     error
}

func (task *Task) ToDomain(uuid string) domain.Task {
	return domain.Task{
		UUID:      uuid,
		Status:    domain.TaskStatus(task.Status),
		CreatedAt: task.CreatedAt,
		StartedAt: task.StartedAt,
		UpdatedAt: task.UpdatedAt,
		Result:    task.Result,
		Error:     task.Error,
	}
}

func FromDomainToTask(task domain.Task) *Task {
	return &Task{
		Status:    string(task.Status),
		CreatedAt: task.CreatedAt,
		StartedAt: task.StartedAt,
		UpdatedAt: task.UpdatedAt,
		Result:    task.Result,
		Error:     task.Error,
	}
}
