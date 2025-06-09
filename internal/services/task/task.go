package task

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type SimulatedTaskService struct {
	log *slog.Logger
}

func NewMockTaskService(log *slog.Logger) *SimulatedTaskService {
	return &SimulatedTaskService{
		log: log,
	}
}

func (m *SimulatedTaskService) Run(ctx context.Context) {
	const op = "MockTaskService.Run"

	log := m.log.With(slog.String("op", op))

	log.InfoContext(ctx, "Running mock task")
	// TOOD: Implement the mock task running logic
	for {
		for _, x := range "-\\|/" {
			fmt.Printf("\r%s", string(x))
			time.Sleep(300 * time.Millisecond)
		}
	}
}
