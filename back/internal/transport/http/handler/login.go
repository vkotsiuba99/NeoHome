package handler

import (
	"encoding/json"
	"io"
	"net/http"
)

func (h *Handler) LoginEmail(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.LoginEmail"
	ctx := r.Context()
	log := h.log.With("op", op)
	log.Info("request started", "method", r.Method, "url", r.URL.String())

	var req LoginEmailReq
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		log.Warn("decode body failed", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		log.Warn("decode body failed", "error", "unexpected trailing payload")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	session, err := h.ReqConv.LoginEmail(ctx, h.authService, req)
	if err != nil {
		handleError(ctx, w, err, log, "login by email failed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(h.RespConv.ToAuth(session)); err != nil {
		log.Error("encode response failed", "error", err.Error())
	}
}
