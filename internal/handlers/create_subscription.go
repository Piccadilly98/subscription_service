package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	dto "github.com/Piccadilly98/subscription_service/internal/models/dto"
	"github.com/Piccadilly98/subscription_service/internal/service"
)

type CreateSubsriptionHandler struct {
	service *service.Service
}

func NewCreateHeandler(s *service.Service) *CreateSubsriptionHandler {
	return &CreateSubsriptionHandler{
		service: s,
	}
}

func (c *CreateSubsriptionHandler) Hander(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(HeaderContentType) != HeaderJson {
		errorResponse(w, fmt.Errorf("invalid header content-type"), http.StatusBadRequest)
		return
	}

	body := &dto.CreateSubscriptionsRequest{}
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		log.Println(err)
		errorResponse(w, fmt.Errorf("invalid body format"), http.StatusBadRequest)
		return
	}

	res, err := c.service.CreateSubsription(r.Context(), body)
	if err != nil {
		log.Println(err)
		errorResponse(w, err, http.StatusBadRequest)
		return
	}

	b, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set(HeaderContentType, HeaderJson)
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}
