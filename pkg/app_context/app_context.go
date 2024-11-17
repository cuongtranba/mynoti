package app_context

import (
	"context"
	"sync"

	"github.com/cuongtranba/mynoti/pkg/logger"
)

type ctxType string

const (
	ctxTypeLogger ctxType = "logger"
)

type AppContext struct {
	context.Context

	logInit sync.Once
}

func New(ctx context.Context) *AppContext {
	return &AppContext{
		Context: ctx,
	}
}

func (c *AppContext) WithContext(ctx context.Context) *AppContext {
	return &AppContext{
		Context: ctx,
	}
}

func (c *AppContext) WithValue(key, value any) *AppContext {
	return &AppContext{
		Context: context.WithValue(c.Context, key, value),
	}
}

func (c *AppContext) WithLogger(logger *logger.Logger) *AppContext {
	return c.WithValue(ctxTypeLogger, logger)
}

func (c *AppContext) Logger() *logger.Logger {
	value, ok := getValue[*logger.Logger](c, ctxTypeLogger)
	if !ok {
		c.logInit.Do(func() {
			logger := logger.NewDefaultLogger()
			c.Context = context.WithValue(c.Context, ctxTypeLogger, logger)
		})
		value, _ = getValue[*logger.Logger](c, ctxTypeLogger)
	}
	return value
}

func getValue[K any](ctx *AppContext, key any) (value K, ok bool) {
	value, ok = ctx.Context.Value(key).(K)
	return
}
