package pool

// Package pool implements a worker pool that manages a set of workers
// to execute tasks concurrently. It provides methods to start the pool,
// submit tasks, and stop the pool gracefully.

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/passwordhash/task-manager-api/internal/domain"
	"github.com/passwordhash/task-manager-api/internal/storage"
	"github.com/passwordhash/task-manager-api/internal/worker"
)

type taskWrapper struct {
	task *domain.Task
	ctx  context.Context
}

type pool struct {
	log       *slog.Logger
	wg        sync.WaitGroup
	workers   int
	taskQueue chan *taskWrapper

	executor    worker.TaskExecutor
	taskStorage storage.Task

	mu         sync.Mutex
	cancelFunc map[string]context.CancelFunc
}

func New(
    log *slog.Logger,
    workers int,
    queueSize int,
    executor worker.TaskExecutor,
    taskStorage storage.Task,
) worker.TaskPool {
	return &pool{
		log:         log,
		workers:     workers,
		taskQueue:   make(chan *taskWrapper, queueSize),
		executor:    executor,
		taskStorage: taskStorage,
		cancelFunc:  make(map[string]context.CancelFunc),
	}
}

func (p *pool) Start(ctx context.Context) {
	const op = "pool.Start"

	log := p.log.With(slog.String("op", op))

	for i := 0; i < p.workers; i++ {
		i := i // but we use go 1.24))
		p.wg.Add(1)
		go p.worker(ctx, i)
	}

	log.Info("Worker pool started", slog.Int("workers", p.workers))
}

func (p *pool) Submit(ctx context.Context, task *domain.Task) error {
	const op = "pool.Submit"

	log := p.log.With(slog.String("op", op), slog.String("task_uuid", task.UUID))

	if task == nil {
		return fmt.Errorf("%s: task cannot be nil", op)
	}

	taskCtx, taskCancel := context.WithCancel(context.Background())
	p.mu.Lock()
	p.cancelFunc[task.UUID] = taskCancel
	p.mu.Unlock()

	tw := &taskWrapper{
		task: task,
		ctx:  taskCtx,
	}

	select {
	case p.taskQueue <- tw:
		log.Debug("Task submitted to the queue")
		return nil
	case <-ctx.Done():
		log.Error("Failed to submit task to the queue", slog.String("error", ctx.Err().Error()))
		return ctx.Err()
	}
}

func (p *pool) Cancel(_ context.Context, taskUUID string) error {
	const op = "pool.Cancel"

	log := p.log.With(slog.String("op", op), slog.String("task_id", taskUUID))

	p.mu.Lock()
	defer p.mu.Unlock()
	cancelFunc, exists := p.cancelFunc[taskUUID]

	if !exists {
		return fmt.Errorf("%s: task %s not found in pool", op, taskUUID)
	}

	cancelFunc()

	log.Debug("Task cancellation requested", slog.String("task_uuid", taskUUID))

	delete(p.cancelFunc, taskUUID)

	return nil
}

func (p *pool) Stop(ctx context.Context) error {
	const op = "pool.Stop"

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

	const op = "pool.worker"

	log := p.log.With(slog.String("op", op), slog.Int("worker_id", id))

	log.Debug("Worker started")

	for {
		select {
		case tw, ok := <-p.taskQueue:
			if !ok {
				log.Debug("Task queue closed, stopping pool")
				return
			}

			wlog := log.With(slog.String("task_uuid", tw.task.UUID))

			wlog.Debug("Received task for execution")

			if err := p.taskStorage.UpdateStatus(ctx, tw.task.UUID, domain.StatusRunning, time.Now()); err != nil {
				// Maybe we should use some retry mechanism here?
				wlog.Error("Failed to update task status to running", slog.String("error", err.Error()))
				continue
			}

			var status domain.TaskStatus
			updatedAt, err := p.executor.Execute(tw.ctx, tw.task)
			if err != nil && errors.Is(err, context.Canceled) {
				wlog.Debug("Task execution cancelled by context")
				status = domain.StatusCancelled
			} else if err != nil && !errors.Is(err, context.Canceled) {
				wlog.Error("Failed to execute task", slog.String("error", err.Error()))
				status = domain.StatusFailed
			} else {
				wlog.Debug("Task executed successfully")
				status = domain.StatusCompleted
			}

			if err := p.taskStorage.UpdateStatus(ctx, tw.task.UUID, status, updatedAt); err != nil {
				// Maybe we should use some retry mechanism here?
				wlog.Error("Failed to update task status after execution", slog.String("error", err.Error()))
				continue
			}
		case <-ctx.Done():
			log.Debug("Worker stopped")
			return
		}
	}
}
