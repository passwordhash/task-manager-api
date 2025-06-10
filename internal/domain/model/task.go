package model

import "time"

type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusRunning              = "running"
	StatusCompleted            = "completed"
	StatusFailed               = "failed"
	StatusCancelled            = "cancelled"
)

type Task struct {
	UUID      string
	Status    TaskStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}
