package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) DeviceLatest(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.DeviceLatest"
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

	latest, err := h.ReqConv.GetLatest(ctx, h.telemetryService, user, deviceID)
	if err != nil {
		handleError(ctx, w, err, log, "get latest failed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(h.RespConv.ToTelemetryList(latest)); err != nil {
		log.Error("encode response failed", "error", err.Error())
	}
}
