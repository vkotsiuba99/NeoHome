package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) Alerts(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.Alerts"
	ctx := r.Context()
	log := h.log.With("op", op)
	log.Info("request started", "method", r.Method, "url", r.URL.String())

	user, err := h.ReqConv.CurrentUser(r)
	if err != nil {
		log.Warn("read current user failed", "error", err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	alerts, err := h.ReqConv.ListAlerts(ctx, h.alertService, r, user)
	if err != nil {
		handleError(ctx, w, err, log, "list alerts failed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(h.RespConv.ToAlerts(alerts)); err != nil {
		log.Error("encode response failed", "error", err.Error())
	}
}
