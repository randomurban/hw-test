package logger

import (
	"log/slog"
	"os"
	"strings"
)

type Logger interface {
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
	TEXT  = "TEXT"
	JSON  = "JSON"
)

func New(level string, logType string) Logger {
	var logLevel slog.Level
	switch strings.ToUpper(level) {
	case DEBUG:
		logLevel = slog.LevelDebug
	case INFO:
		logLevel = slog.LevelInfo
	case WARN:
		logLevel = slog.LevelWarn
	case ERROR:
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}
	logOptions := slog.HandlerOptions{
		AddSource:   false,
		Level:       logLevel,
		ReplaceAttr: nil,
	}
	var logger *slog.Logger

	switch strings.ToUpper(logType) {
	case TEXT:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &logOptions))
	case JSON:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &logOptions))
	default:
		logger = slog.Default() // Simple default logger
	}
	slog.SetDefault(logger)
	return logger
}
