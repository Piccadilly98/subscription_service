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

// GetSummary godoc
// @Summary      Подсчёт суммарной стоимости подписок
// @Description  Возвращает общую стоимость всех подписок за указанный период с возможностью фильтрации.<br><br>Обязательные параметры:<br>• start_date — начало периода (формат "MM-YYYY" или "DD-MM-YYYY")<br>• end_date — конец периода (формат "MM-YYYY" или "DD-MM-YYYY")<br><br>Опциональные параметры:<br>• user_id — фильтр по ID пользователя (UUID)<br>• service_name — фильтр по названию сервиса<br><br>Расчёт включает только подписки, начатые в указанном периоде (учитываются пересечения дат).<br>Если период некорректный (end_date раньше start_date или неверный формат) — возвращается ошибка 400.<br><br>Все ошибки возвращаются в едином формате JSON (ErrorDTO).
// @Tags         subscriptions
// @Produce      json
// @Param        start_date    query  string  true   "Начало периода (MM-YYYY)"      example(08-07-2025)
// @Param        end_date      query  string  true   "Конец периода (MM-YYYY)"       example(12-2025)
// @Param        user_id       query  string  false  "Фильтр по ID пользователя"     example(60601fee-2bf1-4721-ae6f-7636e79a0cba)
// @Param        service_name  query  string  false  "Фильтр по названию сервиса"    example(Yandex Plus)
// @Success      200 {object} dto.SummaryResponse "Суммарная стоимость расчитана и возвращется вместе с заданным периодом"
// @Failure      400 {object} dto.ErrorDTO "Неверные параметры запроса"
// @Failure      500 {object} dto.ErrorDTO "Внутренняя ошибка сервера"
// @Router       /subscriptions/summary [get]
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
