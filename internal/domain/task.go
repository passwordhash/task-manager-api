package domain

import (
	"log/slog"
	"time"
)

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

func (t *Task) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("uuid", t.UUID),
		slog.String("status", string(t.Status)),
		slog.Time("created_at", t.CreatedAt),
		slog.Time("updated_at", t.UpdatedAt),
	)
}
