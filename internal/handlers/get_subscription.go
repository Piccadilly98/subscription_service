package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Piccadilly98/subscription_service/internal/service"
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
	id := chekcURLParam(w, r)
	if id == "" {
		return
	}

	exist, err := g.service.GetExsistBySubID(r.Context(), id)
	if err != nil {
		log.Println(err)
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if !exist {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 page not found"))
		return
	}
	body, err := g.service.GetSubInfoDTOByID(r.Context(), id)
	if err != nil {
		log.Println(err)
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set(HeaderContentType, HeaderJson)
	w.Write(b)
}
