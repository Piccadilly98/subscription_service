package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Piccadilly98/subscription_service/internal/errors_checker"
	"github.com/Piccadilly98/subscription_service/internal/service"
)

type HealthHandler struct {
	serv *service.Service
	ew   *errors_checker.ErrorWorker
}

func NewHealthHandler(serv *service.Service, ew *errors_checker.ErrorWorker) *HealthHandler {
	return &HealthHandler{
		serv: serv,
		ew:   ew,
	}
}

// HealthCheck godoc
// @Summary      Проверка состояния сервиса
// @Description  Возвращает состояние сервера:<br>•Общий статус сервера("ok" или "Service Unavailable")<br>•Статус базы данных в результате пинга("ok" или "does not respond")<br>•Опциональное поле Error с ошибкой от БД<br><br>Эндпоинт предназначен для мониторинга и health-check'ов.
// @Tags         monitoring
// @Produce      json
// @Success      200 {object} dto.CheckHealth "Состояние сервера"
// @Success      503 {object} dto.CheckHealth "Сервис недоступен (проблема с базой данных)"
// @Router       /health-check [get]
func (h *HealthHandler) Handler(w http.ResponseWriter, r *http.Request) {
	status, code := h.serv.CheckHealh(r.Context())

	b, err := json.Marshal(status)
	if err != nil {
		processingError(w, err, h.ew)
		return
	}
	w.Header().Set(HeaderContentType, HeaderJson)
	w.WriteHeader(code)
	w.Write(b)
}
