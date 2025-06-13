package httpapp

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	tasks "github.com/passwordhash/task-manager-api/internal/api/v1/tasks"
	"github.com/passwordhash/task-manager-api/internal/service"
	"github.com/passwordhash/task-manager-api/internal/worker"
)

const shutdownTimeout = 5 * time.Second

type App struct {
	log         *slog.Logger
	taskPool    worker.TaskPool
	taskManager service.TaskService

	port         int
	readTimeout  time.Duration
	writeTimeout time.Duration

	server *http.Server
}

func New(
	log *slog.Logger,
	taskPool worker.TaskPool,
	taskManager service.TaskService,
	port int,
	readTimeout time.Duration,
	writeTimeout time.Duration,
) *App {
	return &App{
		log:          log,
		taskPool:     taskPool,
		taskManager:  taskManager,
		port:         port,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
	}
}

// MustRun starts the HTTP server and panics if it fails to start.
func (a *App) MustRun(ctx context.Context) {
	err := a.Run(ctx)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic("failed to run HTTP server: " + err.Error())
	}
}

// Run starts the HTTP server and listens on the specified port.
func (a *App) Run(ctx context.Context) error {
	const op = "httpapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	log.Info("Starting task pool pool")

	a.taskPool.Start(ctx)

	log.Info("Starting HTTP server")

	router := gin.New()
	router.Use(gin.Recovery())

	api := router.Group("/api")
	v1 := api.Group("/v1")

	tasksHandler := tasks.NewHandler(a.taskManager)

	tasksHandler.RegisterRoutes(v1)

	srv := &http.Server{
		Addr:         ":" + strconv.Itoa(a.port),
		Handler:      router,
		ReadTimeout:  a.readTimeout,
		WriteTimeout: a.writeTimeout,
	}
	a.server = srv

	return srv.ListenAndServe()
}

// Stop gracefully stops the HTTP server.
func (a *App) Stop(ctx context.Context) {
	const op = "httpapp.Stop"

	log := a.log.With(slog.String("op", op))

	log.Info("Stopping HTTP server")

	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	if err := a.taskPool.Stop(ctx); err != nil {
		log.Error("Failed to stop task pool", slog.Any("error", err))
	} else {
		log.Info("Task pool stopped gracefully")
	}

	// Shutdown stops receiving new requests and waits for existing requests to finish.
	if err := a.server.Shutdown(ctx); err != nil {
		log.Error("Failed to gracefully stop HTTP server", slog.Any("error", err))
	} else {
		log.Info("HTTP server stopped gracefully")
	}
}
