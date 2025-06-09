package httpapp

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type App struct {
	log  *slog.Logger
	port int
}

func New(
	_ context.Context,
	log *slog.Logger,
	port int,
) *App {
	return &App{
		log:  log,
		port: port,
	}
}

// MustRun starts the HTTP server and panics if it fails to start.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic("failed to run HTTP server: " + err.Error())
	}
}

// Run starts the HTTP server and listens on the specified port.
func (a *App) Run() error {
	const op = "httpapp.Run"

	log := a.log.With(
		slog.String("op", op), slog.Int("port", a.port))

	log.Info("Starting HTTP server")

	srv, err := a.initHTTPServer()
	if err != nil {
		log.Error("Failed to initialize HTTP server", slog.Any("error", err))
		return err
	}

	return srv.Run(":" + strconv.Itoa(a.port))
}

func (a *App) initHTTPServer() (*gin.Engine, error) {
	router := gin.New()
	router.Use(gin.Recovery())

	for {
		for _, x := range "-\\|/" {
			time.Sleep(300 * time.Millisecond)
			fmt.Printf("\b%s", string(x))
		}
	}

	// router.GET("/example", exampleHandler)

	return router, nil
}

// Stop gracefully stops the HTTP server.
func (a *App) Stop() {
	const op = "httpapp.Stop"

	log := a.log.With(slog.String("op", op))

	log.Info("Stopping HTTP server ...")

	// graceful server shutdown here

	log.Info("HTTP server stopped")
}
