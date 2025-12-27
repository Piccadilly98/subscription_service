package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Piccadilly98/subscription_service/internal/models/dto"
	"github.com/Piccadilly98/subscription_service/internal/service"
)

type UpdateHandler struct {
	serv *service.Service
}

func NewUpdateHandler(serv *service.Service) *UpdateHandler {
	return &UpdateHandler{serv: serv}
}

func (u *UpdateHandler) Handler(w http.ResponseWriter, r *http.Request) {
	if !checkHeaderJson(w, r) {
		return
	}
	id := chekcURLParam(w, r)
	if id == "" {
		return
	}
	exist, err := u.serv.GetExsistBySubID(r.Context(), id)
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
	body := &dto.UpdateSubscriptionRequest{}

	err = json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		log.Println(err)
		errorResponse(w, fmt.Errorf("invalid body format"), http.StatusBadRequest)
		return
	}

	err = u.serv.UpdateSubscription(r.Context(), body, id)
	if err != nil {
		log.Println(err)
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}

	info, err := u.serv.GetSubInfoDTOByID(r.Context(), id)
	if err != nil {
		log.Println(err)
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(info)
	if err != nil {
		log.Println(err)
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set(HeaderContentType, HeaderJson)
	w.Write(b)
}
