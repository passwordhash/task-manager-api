package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

func SetupLogger(env string) *slog.Logger {
	var handler slog.Handler

	w := os.Stdout

	switch env {
	case "dev":
		handler = tint.NewHandler(w, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.TimeOnly,
		})
	case "prod":
		handler = slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level:       slog.LevelInfo,
			ReplaceAttr: replaceTimeFormat,
		})
	default:
		handler = slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelDebug})
	}

	return slog.New(handler)
}

func replaceTimeFormat(_ []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey {
		t := a.Value.Time()
		return slog.String(slog.TimeKey, t.Format(time.RFC3339))
	}
	return a
}
