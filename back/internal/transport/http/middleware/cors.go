package middleware

import (
	"net/http"
	"strings"

	transport "github.com/vkotsiuba99/NeoHome/back/internal/transport/http"
)

func CORS(cfg transport.Config) func(http.Handler) http.Handler {
	allowOrigins := cfg.CORSAllowOrigins
	allowMethods := cfg.CORSAllowMethods
	allowHeaders := cfg.CORSAllowHeaders
	allowCredentials := cfg.CORSAllowCredentials

	return func(next http.Handler) http.Handler {
		if !cfg.CORSEnabled {
			return next
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headers := w.Header()
			headers.Set("Access-Control-Allow-Origin", allowOrigins)
			headers.Set("Access-Control-Allow-Methods", allowMethods)
			headers.Set("Access-Control-Allow-Headers", allowHeaders)
			if allowCredentials {
				headers.Set("Access-Control-Allow-Credentials", "true")
			}

			if r.Method == http.MethodOptions {
				if origin := r.Header.Get("Origin"); len(origin) > 0 && allowOrigins != "*" {
					allowed := false
					for _, candidate := range strings.Split(allowOrigins, ",") {
						if strings.TrimSpace(candidate) == origin {
							allowed = true
							break
						}
					}
					if !allowed {
						w.WriteHeader(http.StatusForbidden)
						return
					}
				}

				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
