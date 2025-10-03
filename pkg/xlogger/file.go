package xlogger

import (
	"log/slog"
	"os"
	"path/filepath"
)

func NewFileHandler(logFile, format string, level slog.Level) (slog.Handler, error) {
	if err := os.MkdirAll(filepath.Dir(logFile), 0755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	var handler slog.Handler
	if format == "json" {
		handler = slog.NewJSONHandler(file, &slog.HandlerOptions{
			Level: level,
		})
	} else {
		handler = slog.NewTextHandler(file, &slog.HandlerOptions{
			Level: level,
		})
	}

	return handler, nil
}
