package handlers

import (
	"net/http"

	"github.com/Piccadilly98/subscription_service/internal/errors_checker"
	"github.com/Piccadilly98/subscription_service/internal/service"
)

type DeleteHanlder struct {
	serv *service.Service
	ew   *errors_checker.ErrorWorker
}

func NewDeleteHandler(serv *service.Service, ew *errors_checker.ErrorWorker) *DeleteHanlder {
	return &DeleteHanlder{
		serv: serv,
		ew:   ew,
	}
}

func (d *DeleteHanlder) Handler(w http.ResponseWriter, r *http.Request) {
	id := checkURLParam(w, r)
	if id == "" {
		return
	}

	err := d.serv.DeleteSubByID(r.Context(), id)
	if err != nil {
		processingError(w, err, d.ew)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
