package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) DeviceThresholdsList(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.DeviceThresholdsList"
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

	thresholds, err := h.ReqConv.GetThresholds(ctx, h.deviceService, user, deviceID)
	if err != nil {
		handleError(ctx, w, err, log, "get thresholds failed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(h.RespConv.ToThresholds(thresholds)); err != nil {
		log.Error("encode response failed", "error", err.Error())
	}
}
