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
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}
func LoggerMiddleware(l *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			ctx := app_context.New(r.Context())
			ctx = ctx.WithContext(ctx.WithLogger(l))
			next.ServeHTTP(w, r.WithContext(ctx))
			l.Info(
				"request",
				logger.String("method", r.Method),
				logger.String("path", r.URL.Path),
				logger.String("status", http.StatusText(lrw.statusCode)),
				logger.Duration("duration", time.Since(startTime)),
			)
		})
	}
}
