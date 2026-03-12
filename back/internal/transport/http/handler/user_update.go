package handler

import (
	"encoding/json"
	"io"
	"net/http"
)

func (h *Handler) UserUpdate(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.UserUpdate"
	ctx := r.Context()
	log := h.log.With("op", op)
	log.Info("request started", "method", r.Method, "url", r.URL.String())

	user, err := h.ReqConv.CurrentUser(r)
	if err != nil {
		log.Warn("read current user failed", "error", err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req UpdateUserReq
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

	updatedUser, err := h.ReqConv.UpdateUser(ctx, h.authService, user, req)
	if err != nil {
		handleError(ctx, w, err, log, "update user failed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(h.RespConv.ToUserWrap(updatedUser)); err != nil {
		log.Error("encode response failed", "error", err.Error())
	}
}
