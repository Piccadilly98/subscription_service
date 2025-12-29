package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Piccadilly98/subscription_service/internal/errors_checker"
	"github.com/Piccadilly98/subscription_service/internal/models/dto"
	"github.com/Piccadilly98/subscription_service/internal/service"
)

type GetSummaryHandler struct {
	serv *service.Service
	ew   *errors_checker.ErrorWorker
}

func NewGetSummaryHandler(serv *service.Service, ew *errors_checker.ErrorWorker) *GetSummaryHandler {
	return &GetSummaryHandler{
		serv: serv,
		ew:   ew,
	}
}

func (g *GetSummaryHandler) Handler(w http.ResponseWriter, r *http.Request) {
	params := g.checkURLParams(r)

	res, err := g.serv.GetSummarySubs(r.Context(), params)
	if err != nil {
		processingError(w, err, g.ew)
		return
	}

	b, err := json.Marshal(res)
	if err != nil {
		processingError(w, err, g.ew)
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
