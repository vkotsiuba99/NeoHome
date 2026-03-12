package middleware

import (
	"log/slog"
	"net/http"
	"time"

	transport "github.com/vkotsiuba99/NeoHome/back/internal/transport/http"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
	bytes      int
}

func (w *statusRecorder) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusRecorder) Write(data []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}

	size, err := w.ResponseWriter.Write(data)
	w.bytes += size
	return size, err
}

func AccessLog(cfg transport.Config, log *slog.Logger) func(http.Handler) http.Handler {
	if !cfg.LoggerEnabled {
		return func(next http.Handler) http.Handler { return next }
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			wrappedWriter := &statusRecorder{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(wrappedWriter, r)

			log.Info("request processed",
				"method", r.Method,
				"path", r.URL.Path,
				"query", r.URL.RawQuery,
				"status_code", wrappedWriter.statusCode,
				"duration_ms", time.Since(startTime).Milliseconds(),
				"bytes", wrappedWriter.bytes,
				"user_agent", r.UserAgent(),
				"remote_addr", r.RemoteAddr,
			)
		})
	}
}
