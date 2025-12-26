package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Piccadilly98/subscription_service/internal/models/dto"
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

//checker errors and get code and err
