package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/vkotsiuba99/NeoHome/back/internal/service"
)

func handleError(_ context.Context, w http.ResponseWriter, err error, log *slog.Logger, defaultMsg string) {
	switch {
	case errors.Is(err, ErrValidation), errors.Is(err, service.ErrValidation):
		log.Warn(defaultMsg, "error", err.Error(), "status_code", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, ErrUnauthorized), errors.Is(err, service.ErrUnauthorized):
		log.Warn(defaultMsg, "error", err.Error(), "status_code", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
	case errors.Is(err, ErrForbidden), errors.Is(err, service.ErrForbidden):
		log.Warn(defaultMsg, "error", err.Error(), "status_code", http.StatusForbidden)
		w.WriteHeader(http.StatusForbidden)
	case errors.Is(err, ErrNotFound), errors.Is(err, service.ErrNotFound):
		log.Warn(defaultMsg, "error", err.Error(), "status_code", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
	case errors.Is(err, ErrConflict), errors.Is(err, service.ErrConflict):
		log.Warn(defaultMsg, "error", err.Error(), "status_code", http.StatusConflict)
		w.WriteHeader(http.StatusConflict)
	default:
		log.Error(defaultMsg, "error", err.Error(), "status_code", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
