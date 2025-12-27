package handlers

import (
	"encoding/json"
	"net/http"

	error_worker "github.com/Piccadilly98/subscription_service/internal/errorWorker"
	"github.com/Piccadilly98/subscription_service/internal/service"
)

type HealthHandler struct {
	serv      *service.Service
	errWorker *error_worker.ErrorWorker
}

func NewHealthHandler(serv *service.Service, errWorker *error_worker.ErrorWorker) *HealthHandler {
	return &HealthHandler{
		serv:      serv,
		errWorker: errWorker,
	}
}

func (h *HealthHandler) Handler(w http.ResponseWriter, r *http.Request) {
	status := h.serv.CheckHealh(r.Context())

	b, err := json.Marshal(status)
	if err != nil {
		formatingErrorAndWriteError(h.errWorker, w, err, http.StatusOK, ErrorInvalidRequest)
		return
	}
	w.Header().Set(HeaderContentType, HeaderJson)
	w.Write(b)
}
