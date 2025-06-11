package model

import (
	"time"

	"github.com/passwordhash/task-manager-api/internal/domain"
)

type Task struct {
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (task *Task) ToDomain() domain.Task {
	return domain.Task{
		Status:    domain.TaskStatus(task.Status),
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
}

func FromDomainToTask(task domain.Task) *Task {
	return &Task{
		Status:    string(task.Status),
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
}
