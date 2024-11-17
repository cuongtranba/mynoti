package middleware

import (
	"net/http"

	"github.com/cuongtranba/mynoti/pkg/app_context"
	"github.com/cuongtranba/mynoti/pkg/logger"
)

func LoggerMiddleware(logger *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add the logger to the context
			ctx := app_context.New(r.Context())
			ctx = ctx.WithContext(ctx.WithLogger(logger))
			// Pass the request with the new context to the next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
