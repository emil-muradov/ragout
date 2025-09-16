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