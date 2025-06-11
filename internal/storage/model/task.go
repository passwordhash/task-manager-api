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

func FromDomainToTask(task domain.Task) *Task {
	return &Task{
		Status:    string(task.Status),
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
}
