package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) AlertResolve(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.AlertResolve"
	ctx := r.Context()
	log := h.log.With("op", op)
	log.Info("request started", "method", r.Method, "url", r.URL.String())

	user, err := h.ReqConv.CurrentUser(r)
	if err != nil {
		log.Warn("read current user failed", "error", err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	alertID, err := h.ReqConv.PathInt64(r, "alertId")
	if err != nil {
		log.Warn("parse alert id failed", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	alert, err := h.ReqConv.ResolveAlert(ctx, h.alertService, user, alertID)
	if err != nil {
		handleError(ctx, w, err, log, "resolve alert failed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(h.RespConv.ToAlertWrap(alert)); err != nil {
		log.Error("encode response failed", "error", err.Error())
	}
}
