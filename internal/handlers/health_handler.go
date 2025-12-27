package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Piccadilly98/subscription_service/internal/service"
)

type HealthHandler struct {
	serv *service.Service
}

func NewHealthHandler(serv *service.Service) *HealthHandler {
	return &HealthHandler{serv: serv}
}

func (h *HealthHandler) Handler(w http.ResponseWriter, r *http.Request) {
	status := h.serv.CheckHealh(r.Context())

	b, err := json.Marshal(status)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set(HeaderContentType, HeaderJson)
	w.Write(b)
}
