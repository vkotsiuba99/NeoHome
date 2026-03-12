package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	transport "github.com/vkotsiuba99/NeoHome/back/internal/transport/http"
)

func Recover(cfg transport.Config, log *slog.Logger) func(http.Handler) http.Handler {
	if !cfg.RecoverEnabled {
		return func(next http.Handler) http.Handler { return next }
	}

	recoverMessage := cfg.RecoverMessage
	includeStack := cfg.RecoverStack

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					fields := []any{
						"method", r.Method,
						"path", r.URL.Path,
						"panic", fmt.Sprint(rec),
					}
					if includeStack {
						fields = append(fields, "stack", string(debug.Stack()))
					}

					log.Error("panic recovered", fields...)
					http.Error(w, recoverMessage, http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
