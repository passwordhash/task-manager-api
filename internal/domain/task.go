package domain

import (
	"fmt"
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
	StartedAt time.Time
	UpdatedAt time.Time
}

func (t *Task) RunningDuration() time.Duration {
	fmt.Println(t.StartedAt.String())
	switch t.Status {
	case StatusPending:
		return 0
	case StatusRunning:
		return time.Since(t.StartedAt)
	case StatusCompleted, StatusFailed, StatusCancelled:
		return t.UpdatedAt.Sub(t.StartedAt)
	default:
		return 0
	}
}

func (t *Task) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("uuid", t.UUID),
		slog.String("status", string(t.Status)),
		slog.Time("created_at", t.CreatedAt),
		slog.Time("updated_at", t.UpdatedAt),
	)
}
