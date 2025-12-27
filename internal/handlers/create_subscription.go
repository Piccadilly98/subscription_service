package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	error_worker "github.com/Piccadilly98/subscription_service/internal/errorWorker"
	dto "github.com/Piccadilly98/subscription_service/internal/models/dto"
	"github.com/Piccadilly98/subscription_service/internal/service"
)

type CreateSubsriptionHandler struct {
	service   *service.Service
	errWorker *error_worker.ErrorWorker
}

func NewCreateHandler(s *service.Service, errWorker *error_worker.ErrorWorker) *CreateSubsriptionHandler {
	return &CreateSubsriptionHandler{
		service:   s,
		errWorker: errWorker,
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
		formatingErrorAndWriteError(c.errWorker, w, err, http.StatusCreated, ErrorInvalidBody)
		return
	}

	b, err := json.Marshal(res)
	if err != nil {
		formatingErrorAndWriteError(c.errWorker, w, err, http.StatusCreated, "")
		return
	}

	w.Header().Set(HeaderContentType, HeaderJson)
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}
