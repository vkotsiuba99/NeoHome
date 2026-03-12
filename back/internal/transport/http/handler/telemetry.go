package handler

import (
	"encoding/json"
	"io"
	"net/http"
)

func (h *Handler) Telemetry(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.Telemetry"
	ctx := r.Context()
	log := h.log.With("op", op)
	log.Info("request started", "method", r.Method, "url", r.URL.String())

	var body telemetryPayloadBody
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		log.Warn("decode body failed", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		log.Warn("decode body failed", "error", "unexpected trailing payload")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req := h.ReqConv.ToTelemetryIngest(body)

	result, err := h.ReqConv.IngestTelemetry(ctx, h.telemetryService, req)
	if err != nil {
		handleError(ctx, w, err, log, "ingest telemetry failed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(h.RespConv.ToTelemetryIngest(result)); err != nil {
		log.Error("encode response failed", "error", err.Error())
	}
}
