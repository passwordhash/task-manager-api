package pool

// Package pool implements a worker pool that manages a set of workers
// to execute tasks concurrently. It provides methods to start the pool,
// submit tasks, and stop the pool gracefully.

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/passwordhash/task-manager-api/internal/domain"
	"github.com/passwordhash/task-manager-api/internal/worker"
)

type pool struct {
	log       *slog.Logger
	workers   int
	taskQueue chan *domain.Task
	wg        sync.WaitGroup
	executor  worker.TaskExecutor
}

func New(
	log *slog.Logger,
	queueSize,
	workers int,
	executor worker.TaskExecutor,
) worker.TaskPool {
	return &pool{
		log:       log,
		workers:   workers,
		taskQueue: make(chan *domain.Task, queueSize),
		executor:  executor,
	}
}

func (p *pool) Start(ctx context.Context) {
	const op = "workerpool.Start"

	log := p.log.With(slog.String("op", op))

	for i := 0; i < p.workers; i++ {
		i := i // but we use go 1.24))
		p.wg.Add(1)
		go p.worker(ctx, i)
	}

	log.Info("Worker pool started", slog.Int("workers", p.workers))
}

func (p *pool) Submit(ctx context.Context, task *domain.Task) error {
	const op = "workerpool.Submit"

	log := p.log.With(slog.String("op", op), slog.String("task_uuid", task.UUID))

	select {
	case p.taskQueue <- task:
		log.Debug("Task submitted to the queue")
		return nil
	case <-ctx.Done():
		log.Error("Failed to submit task to the queue", slog.String("error", ctx.Err().Error()))
		return ctx.Err()
	}
}

func (p *pool) Stop(ctx context.Context) error {
	const op = "workerpool.Stop"

	log := p.log.With(slog.String("op", op))

	close(p.taskQueue)

	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		log.Warn("Context cancelled before workers finished")
		return fmt.Errorf("%s: context cancelled before workers finished: %w", op, ctx.Err())
	case <-done:
		log.Info("Worker pool stopped")
		return nil
	}
}

func (p *pool) worker(ctx context.Context, id int) {
	defer p.wg.Done()

	const op = "workerpool.pool"

	log := p.log.With(slog.String("op", op), slog.Int("worker_id", id))

	log.Debug("Worker started")

	for {
		select {
		case task, ok := <-p.taskQueue:
			if !ok {
				log.Debug("Task queue closed, stopping pool")
				return
			}

			p.executor.Execute(ctx, task)
		case <-ctx.Done():
			log.Debug("Worker stopped")
			return
		}
	}
}
