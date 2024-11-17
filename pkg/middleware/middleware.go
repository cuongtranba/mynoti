package middleware

import (
	"net/http"
	"time"

	"github.com/cuongtranba/mynoti/pkg/app_context"
	"github.com/cuongtranba/mynoti/pkg/logger"
)

func ContextMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := app_context.New(r.Context())
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	req        *http.Request
	statusCode int
	logger     *logger.Logger
	startTime  time.Time
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	defer func() {
		r := lrw.req
		startTime := lrw.startTime
		l := lrw.logger
		l.Info(
			"request",
			logger.String("method", r.Method),
			logger.String("path", r.URL.Path),
			logger.String("status", http.StatusText(lrw.statusCode)),
			logger.Duration("duration", time.Since(startTime)),
		)
	}()
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}
func LoggerMiddleware(l *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lrw := &loggingResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
				logger:         l,
				req:            r,
				startTime:      time.Now(),
			}
			ctx := app_context.New(r.Context())
			ctx = ctx.WithContext(ctx.WithLogger(l))
			next.ServeHTTP(lrw, r.WithContext(ctx))
		})
	}
}
