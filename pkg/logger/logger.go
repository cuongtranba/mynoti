package logger

import (
	"log/slog"
	"os"
	"sync"
)

var (
	logger *Logger
	doOnce sync.Once
)

type Logger struct {
	*slog.Logger
}

func NewLogger(appName string) *Logger {
	v := &Logger{
		Logger: slog.With(
			slog.String("app", appName)),
	}
	return v
}

func NewDefaultLogger() *Logger {
	doOnce.Do(func() {
		logger = NewLogger(os.Getenv("APP_NAME"))
	})
	return logger
}
