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
	loggerMap sync.Map
)

type Logger struct {
	*slog.Logger
}

func (l *Logger) Fatal(err error) {
	l.Error(err.Error())
	os.Exit(1)
}

func NewLogger(appName string) *Logger {
	logger := &Logger{
		Logger: slog.With(
			slog.String("app", appName)),
	}
	l, _ := loggerMap.LoadOrStore(appName, logger)
	return l.(*Logger)
}

func NewDefaultLogger() *Logger {
	logger := NewLogger("app")
	return logger
}
