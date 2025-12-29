package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Piccadilly98/subscription_service/internal/errors_checker"
	"github.com/Piccadilly98/subscription_service/internal/models/dto"
	"github.com/Piccadilly98/subscription_service/internal/service"
)

type UpdateHandler struct {
	serv *service.Service
	ew   *errors_checker.ErrorWorker
}

func NewUpdateHandler(serv *service.Service, ew *errors_checker.ErrorWorker) *UpdateHandler {
	return &UpdateHandler{
		serv: serv,
		ew:   ew,
	}
}

func (u *UpdateHandler) Handler(w http.ResponseWriter, r *http.Request) {
	if !checkHeaderJson(w, r) {
		return
	}
	id := checkURLParam(w, r)
	if id == "" {
		return
	}
	exist, err := u.serv.GetExsistBySubID(r.Context(), id)
	if err != nil {
		processingError(w, err, u.ew)
		return
	}
	if !exist {
		processingError(w, errors.New("subscription not found"), u.ew)
		return
	}
	body := &dto.UpdateSubscriptionRequest{}

	err = json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		log.Println(err)
		errorResponse(w, fmt.Errorf("invalid body format"), http.StatusBadRequest)
		return
	}

	res, err := u.serv.UpdateSubscription(r.Context(), body, id)
	if err != nil {
		processingError(w, err, u.ew)
		return
	}

	b, err := json.Marshal(res)
	if err != nil {
		processingError(w, err, u.ew)
		return
	}
	w.Header().Set(HeaderContentType, HeaderJson)
	w.Write(b)
}
