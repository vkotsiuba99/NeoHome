package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/vkotsiuba99/NeoHome/back/pkg/jwt"
)

const userHeaderKey = "X-User"

func Auth(secret []byte, log *slog.Logger) func(http.Handler) http.Handler {
	log = log.With(
		"layer", "transport",
		"component", "auth_middleware",
		"op", "middleware.Auth",
	)

	const authorizationHeaderKey = "Authorization"
	const bearerScheme = "Bearer"

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Header.Del(userHeaderKey)

			if len(secret) == 0 {
				log.Error("internal error",
					"reason", "empty_jwt_secret",
					"method", r.Method,
					"path", r.URL.Path,
				)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			authHeaderValue := r.Header.Get(authorizationHeaderKey)
			if len(authHeaderValue) == 0 {
				log.Info("unauthorized",
					"reason", "missing_authorization_header",
					"method", r.Method,
					"path", r.URL.Path,
				)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			parts := strings.Fields(authHeaderValue)
			if len(parts) != 2 || !strings.EqualFold(parts[0], bearerScheme) || len(parts[1]) == 0 {
				log.Info("unauthorized",
					"reason", "invalid_authorization_header_format",
					"method", r.Method,
					"path", r.URL.Path,
				)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			claims, parseErr := jwt.ParseToken(parts[1], secret)
			if parseErr != nil {
				log.Info("unauthorized",
					"reason", "invalid_token",
					"method", r.Method,
					"path", r.URL.Path,
				)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			userData, claimsErr := jwt.ParseUserClaims(claims)
			if claimsErr != nil {
				log.Info("unauthorized",
					"reason", "invalid_claims",
					"method", r.Method,
					"path", r.URL.Path,
				)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			rawUserData, marshalErr := json.Marshal(userData)
			if marshalErr != nil {
				log.Error("internal error",
					"reason", "marshal_user_failed",
					"method", r.Method,
					"path", r.URL.Path,
					"error", marshalErr.Error(),
				)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Header.Set(userHeaderKey, string(rawUserData))
			next.ServeHTTP(w, r)
		})
	}
}
