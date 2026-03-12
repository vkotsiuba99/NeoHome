package handler

import (
	"encoding/json"
	"io"
	"net/http"
)

func (h *Handler) DeviceUpdate(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.DeviceUpdate"
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

	var raw deviceUpdateRawBody
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&raw); err != nil {
		log.Warn("decode body failed", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		log.Warn("decode body failed", "error", "unexpected trailing payload")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req := h.ReqConv.ToDeviceUpdate(raw)

	device, err := h.ReqConv.UpdateDevice(ctx, h.deviceService, user, deviceID, req)
	if err != nil {
		handleError(ctx, w, err, log, "update device failed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(h.RespConv.ToDeviceWrap(device)); err != nil {
		log.Error("encode response failed", "error", err.Error())
	}
}
