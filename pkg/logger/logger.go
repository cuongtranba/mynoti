package logger

import (
    "log"
)

// Info logs an info message
func Info(message string) {
    log.Println("[INFO]", message)
}
