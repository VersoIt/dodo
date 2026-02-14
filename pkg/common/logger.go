package common

import (
	"log/slog"
	"os"
)

func NewLogger(serviceName string) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		AddSource: true,
	}).WithAttrs([]slog.Attr{
		slog.String("service", serviceName),
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}
