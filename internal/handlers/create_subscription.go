package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Piccadilly98/subscription_service/internal/errors_checker"
	dto "github.com/Piccadilly98/subscription_service/internal/models/dto"
	"github.com/Piccadilly98/subscription_service/internal/service"
)

type CreateSubsriptionHandler struct {
	service *service.Service
	ew      *errors_checker.ErrorWorker
}

func NewCreateHandler(s *service.Service, ew *errors_checker.ErrorWorker) *CreateSubsriptionHandler {
	return &CreateSubsriptionHandler{
		service: s,
		ew:      ew,
	}
}

func (c *CreateSubsriptionHandler) Handler(w http.ResponseWriter, r *http.Request) {
	if !checkHeaderJson(w, r) {
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
		processingError(w, err, c.ew)
		return
	}

	b, err := json.Marshal(res)
	if err != nil {
		processingError(w, err, c.ew)
		return
	}

	w.Header().Set(HeaderContentType, HeaderJson)
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}
