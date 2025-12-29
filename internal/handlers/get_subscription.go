package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Piccadilly98/subscription_service/internal/errors_checker"
	"github.com/Piccadilly98/subscription_service/internal/service"
)

type GetSubscriptionHandler struct {
	service *service.Service
	ew      *errors_checker.ErrorWorker
}

func NewGetHandler(s *service.Service, ew *errors_checker.ErrorWorker) *GetSubscriptionHandler {
	return &GetSubscriptionHandler{
		service: s,
		ew:      ew,
	}
}

func (g *GetSubscriptionHandler) Handler(w http.ResponseWriter, r *http.Request) {
	id := checkURLParam(w, r)
	if id == "" {
		return
	}

	exist, err := g.service.GetExsistBySubID(r.Context(), id)
	if err != nil {
		processingError(w, err, g.ew)
		return
	}

	if !exist {
		processingError(w, errors.New("subscription not found"), g.ew)
		return
	}
	body, err := g.service.GetSubInfoDTOByID(r.Context(), id)
	if err != nil {
		processingError(w, err, g.ew)
		return
	}

	b, err := json.Marshal(body)
	if err != nil {
		processingError(w, err, g.ew)
		return
	}
	w.Header().Set(HeaderContentType, HeaderJson)
	w.Write(b)
}
