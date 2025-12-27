package handlers

import (
	"net/http"

	"github.com/Piccadilly98/subscription_service/internal/service"
)

type DeleteHanlder struct {
	serv *service.Service
}

func NewDeleteHandler(serv *service.Service) *DeleteHanlder {
	return &DeleteHanlder{serv: serv}
}

func (d *DeleteHanlder) Handler(w http.ResponseWriter, r *http.Request) {
	id := chekcURLParam(w, r)
	if id == "" {
		return
	}

	err := d.serv.DeleteSubByID(r.Context(), id)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
