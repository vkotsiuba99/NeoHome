package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) TelemetryDevice(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.TelemetryDevice"
	ctx := r.Context()
	log := h.log.With("op", op)
	log.Info("request started", "method", r.Method, "url", r.URL.String())

	user, err := h.ReqConv.CurrentUser(r)
	if err != nil {
		log.Warn("read current user failed", "error", err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	deviceID, err := h.ReqConv.PathInt64(r, "deviceId")
	if err != nil {
		log.Warn("parse device id failed", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	telemetryRows, err := h.ReqConv.ListTelemetry(ctx, h.telemetryService, user, deviceID, r)
	if err != nil {
		handleError(ctx, w, err, log, "list telemetry failed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(h.RespConv.ToTelemetryList(telemetryRows)); err != nil {
		log.Error("encode response failed", "error", err.Error())
	}
}
