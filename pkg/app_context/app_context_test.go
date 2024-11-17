package app_context

import (
	"context"
	"testing"

	"github.com/cuongtranba/mynoti/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestDefaultLoggerInitialization(t *testing.T) {
	// Arrange
	ctx := context.Background()
	appCtx := New(ctx)

	// Act
	logger1 := appCtx.Logger()
	logger2 := appCtx.Logger()

	// Assert
	assert.NotNil(t, logger1, "Logger should not be nil")
	assert.Equal(t, logger1, logger2, "Logger should be initialized only once")
}

func TestCustomLoggerSetAndRetrieve(t *testing.T) {
	// Arrange
	ctx := context.Background()
	appCtx := New(ctx)
	customLogger := logger.NewDefaultLogger() // Replace with appropriate method to create a custom logger

	// Act
	appCtx.WithLogger(customLogger)
	retrievedLogger := appCtx.Logger()

	// Assert
	assert.Equal(t, customLogger, retrievedLogger, "Custom logger should be retrievable after being set")
}

func TestDefaultLoggerDoesNotReinitialize(t *testing.T) {
	// Arrange
	ctx := context.Background()
	appCtx := New(ctx)

	// Act
	logger1 := appCtx.Logger()
	logger2 := appCtx.Logger()

	// Assert
	assert.Same(t, logger1, logger2, "Default logger should not be reinitialized")
}

func TestThreadSafetyOfLoggerInitialization(t *testing.T) {
	// Arrange
	ctx := context.Background()
	appCtx := New(ctx)
	const goroutines = 100
	results := make(chan *logger.Logger, goroutines)

	// Act
	for i := 0; i < goroutines; i++ {
		go func() {
			results <- appCtx.Logger()
		}()
	}

	// Collect all results
	var firstLogger *logger.Logger
	allSame := true
	for i := 0; i < goroutines; i++ {
		l := <-results
		if firstLogger == nil {
			firstLogger = l
		} else if firstLogger != l {
			allSame = false
		}
	}

	// Assert
	assert.True(t, allSame, "All goroutines should receive the same logger instance")
	assert.NotNil(t, firstLogger, "Logger should not be nil")
}
