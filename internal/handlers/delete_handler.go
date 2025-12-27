package handlers

import (
	"net/http"

	error_worker "github.com/Piccadilly98/subscription_service/internal/errorWorker"
	"github.com/Piccadilly98/subscription_service/internal/service"
)

type DeleteHanlder struct {
	serv      *service.Service
	errWorker *error_worker.ErrorWorker
}

func NewDeleteHandler(serv *service.Service, errWorker *error_worker.ErrorWorker) *DeleteHanlder {
	return &DeleteHanlder{
		serv:      serv,
		errWorker: errWorker,
	}
}

func (d *DeleteHanlder) Handler(w http.ResponseWriter, r *http.Request) {
	id := checkURLParam(w, r)
	if id == "" {
		return
	}

	err := d.serv.DeleteSubByID(r.Context(), id)
	if err != nil {
		formatingErrorAndWriteError(d.errWorker, w, err, http.StatusNoContent, "")
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
