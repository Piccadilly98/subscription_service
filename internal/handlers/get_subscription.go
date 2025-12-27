package handlers

import (
	"encoding/json"
	"net/http"

	error_worker "github.com/Piccadilly98/subscription_service/internal/errorWorker"
	"github.com/Piccadilly98/subscription_service/internal/service"
)

type GetSubscriptionHandler struct {
	service   *service.Service
	errWorker *error_worker.ErrorWorker
}

func NewGetHandler(s *service.Service, errWorker *error_worker.ErrorWorker) *GetSubscriptionHandler {
	return &GetSubscriptionHandler{
		service:   s,
		errWorker: errWorker,
	}
}

func (g *GetSubscriptionHandler) Handler(w http.ResponseWriter, r *http.Request) {
	id := checkURLParam(w, r)
	if id == "" {
		return
	}

	exist, err := g.service.GetExsistBySubID(r.Context(), id)
	if err != nil {
		formatingErrorAndWriteError(g.errWorker, w, err, http.StatusOK, ErrorInvalidBody)
		return
	}

	if !exist {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 page not found"))
		return
	}
	body, err := g.service.GetSubInfoDTOByID(r.Context(), id)
	if err != nil {
		formatingErrorAndWriteError(g.errWorker, w, err, http.StatusOK, ErrorInvalidBody)
		return
	}

	b, err := json.Marshal(body)
	if err != nil {
		formatingErrorAndWriteError(g.errWorker, w, err, http.StatusOK, ErrorInvalidBody)
		return
	}
	w.Header().Set(HeaderContentType, HeaderJson)
	w.Write(b)
}
