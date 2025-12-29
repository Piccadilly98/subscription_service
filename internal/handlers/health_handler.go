package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Piccadilly98/subscription_service/internal/errors_checker"
	"github.com/Piccadilly98/subscription_service/internal/service"
)

type HealthHandler struct {
	serv *service.Service
	ew   *errors_checker.ErrorWorker
}

func NewHealthHandler(serv *service.Service, ew *errors_checker.ErrorWorker) *HealthHandler {
	return &HealthHandler{
		serv: serv,
		ew:   ew,
	}
}

func (h *HealthHandler) Handler(w http.ResponseWriter, r *http.Request) {
	status := h.serv.CheckHealh(r.Context())

	b, err := json.Marshal(status)
	if err != nil {
		processingError(w, err, h.ew)
		return
	}
	w.Header().Set(HeaderContentType, HeaderJson)
	w.Write(b)
}
