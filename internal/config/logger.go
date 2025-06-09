package config

import (
	"log/slog"
	"os"
)

func SetupLogger(env string) *slog.Logger {
	var level slog.Level
	switch env {
	case "dev":
		level = slog.LevelDebug
	case "prod":
		level = slog.LevelInfo
	default:
		level = slog.LevelInfo
	}

	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
}
