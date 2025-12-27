package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Piccadilly98/subscription_service/internal/models/dto"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func errorResponse(w http.ResponseWriter, err error, code int) {
	w.Header().Set(HeaderContentType, HeaderJson)
	w.WriteHeader(code)
	resp := dto.NewErrorDto(err)
	b, err := json.Marshal(resp)
	if err != nil {
		return
	}
	w.Write(b)
}

func checkHeaderJson(w http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get(HeaderContentType) != HeaderJson {
		errorResponse(w, fmt.Errorf("invalid header content-type"), http.StatusBadRequest)
		return false
	}
	return true
}

func chekcURLParam(w http.ResponseWriter, r *http.Request) string {
	id := chi.URLParam(r, URLParam)
	if id == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 page not found"))
		return ""
	}
	if _, err := uuid.Parse(id); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 page not found"))
		return ""
	}
	return id
}
