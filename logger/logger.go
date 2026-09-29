package logger

import (
	"log/slog"
	"os"
)

func LoadLogger() {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	slog.SetDefault(logger)
}
