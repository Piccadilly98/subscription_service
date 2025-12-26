package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Piccadilly98/subscription_service/internal/service"
	"github.com/go-chi/chi/v5"
)

type GetSubscriptionHandler struct {
	service *service.Service
}

func NewGetHandler(s *service.Service) *GetSubscriptionHandler {
	return &GetSubscriptionHandler{
		service: s,
	}
}

func (g *GetSubscriptionHandler) Handler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, URLParam)
	if id == "" {
		errorResponse(w, fmt.Errorf("invalid id"), http.StatusBadRequest)
		return
	}

	exist, err := g.service.GetExsistBySubID(r.Context(), id)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if !exist {
		errorResponse(w, fmt.Errorf("invalid id"), http.StatusBadRequest)
		return
	}
	body, err := g.service.GetSubInfoDTOByID(r.Context(), id)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(body)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set(HeaderContentType, HeaderJson)
	w.Write(b)
}
