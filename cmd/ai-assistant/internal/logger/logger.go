package logger

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func CreateLogger() *slog.Logger {
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return Logger
}

func Error(msg string, args ...any) {
	Logger.Error(msg, args...)
}