package app

import (
	"context"
	"log/slog"

	httpapp "github.com/passwordhash/task-manager-api/internal/app/http"
	"github.com/passwordhash/task-manager-api/internal/config"
)

type App struct {
	HTTPSrv *httpapp.App
}

func New(
	ctx context.Context,
	log *slog.Logger,
	cfg *config.Config,
) *App {
	httpApp := httpapp.New(
		ctx,
		log,
		cfg.App.Port,
	)

	return &App{
		HTTPSrv: httpApp,
	}
}
