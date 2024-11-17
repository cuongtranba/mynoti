package logger

import (
	"log/slog"
	"os"
	"sync"
)

var (
	String       = slog.String
	Bool         = slog.Bool
	Int          = slog.Int
	Int64        = slog.Int64
	Uint64       = slog.Uint64
	Float64      = slog.Float64
	Time         = slog.Time
	Any          = slog.Any
	BoolValue    = slog.BoolValue
	Float64Value = slog.Float64Value
	Duration     = slog.Duration
)

var (
	logger *Logger
	doOnce sync.Once
)

type Logger struct {
	*slog.Logger
}

func (l *Logger) Fatal(err error) {
	l.Error(err.Error())
	os.Exit(1)
}

func NewLogger(appName string) *Logger {
	doOnce.Do(func() {
		logger = &Logger{
			Logger: slog.With(
				slog.String("app", appName)),
		}
	})
	return logger
}

func NewDefaultLogger() *Logger {
	logger = NewLogger("app")
	return logger
}
