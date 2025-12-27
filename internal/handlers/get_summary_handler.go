package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Piccadilly98/subscription_service/internal/models/dto"
	"github.com/Piccadilly98/subscription_service/internal/service"
)

type GetSummaryHandler struct {
	serv *service.Service
}

func NewGetSummaryHandler(serv *service.Service) *GetSummaryHandler {
	return &GetSummaryHandler{serv: serv}
}

func (g *GetSummaryHandler) Handler(w http.ResponseWriter, r *http.Request) {
	params := g.checkURLParams(r)

	res, err := g.serv.GetSummarySubs(r.Context(), params)
	if err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}

	b, err := json.Marshal(res)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set(HeaderContentType, HeaderJson)
	w.Write(b)
}

func (g *GetSummaryHandler) checkURLParams(r *http.Request) *dto.QueryParamsSummary {
	model := &dto.QueryParamsSummary{}
	if str := r.URL.Query().Get(QueryUserID); str != "" {
		model.UserID = &str
	}
	if str := r.URL.Query().Get(QueryServiceName); str != "" {
		model.ServiceName = &str
	}
	if str := r.URL.Query().Get(QueryStartDate); str != "" {
		model.StartDate = str
	}
	if str := r.URL.Query().Get(QueryEndDate); str != "" {
		model.EndDate = str
	}

	return model
}
