package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Piccadilly98/subscription_service/internal/errors_checker"
	"github.com/Piccadilly98/subscription_service/internal/service"
)

type GetSubscriptionHandler struct {
	service *service.Service
	ew      *errors_checker.ErrorWorker
}

func NewGetHandler(s *service.Service, ew *errors_checker.ErrorWorker) *GetSubscriptionHandler {
	return &GetSubscriptionHandler{
		service: s,
		ew:      ew,
	}
}

// GetSubscription godoc
// @Summary      Получить данные подписки
// @Description  Получает данные по ID подписки(UUID).<br><br>Сервис выполняет проверки:<br>• ID не пустой и является валидным UUID<br>• Подписка с таким ID существует<br><br>При успешном запросе возвращает данные подписки с рассчитанным статусом "active", "ended" или "not started"<br><br>Все ошибки возвращаются в формате JSON(dto.ErrorDTO)
//
// @Tags         subscriptions
// @Produce      json
// @Param        id  path  string  true  "ID подписки (UUID)"  example(8b2d4644-5688-47f7-a8b6-7e6d4f55a951)
// @Success      200 {object} dto.SubscriptionResponse "Подписка успешно получена"
// @Failure      400 {object} dto.ErrorDTO "ID не uuid или пустой"
// @Failure      404 {object} dto.ErrorDTO "ID не существует"
// @Failure      500 {object} dto.ErrorDTO "Ошибка сервера"
// @Router       /subscriptions/{id} [get]
func (g *GetSubscriptionHandler) Handler(w http.ResponseWriter, r *http.Request) {
	id := checkURLParam(w, r, g.ew)
	if id == "" {
		return
	}

	exist, err := g.service.GetExsistBySubID(r.Context(), id)
	if err != nil {
		processingError(w, err, g.ew)
		return
	}

	if !exist {
		processingError(w, errors.New("subscription not found"), g.ew)
		return
	}
	body, err := g.service.GetSubInfoDTOByID(r.Context(), id)
	if err != nil {
		processingError(w, err, g.ew)
		return
	}

	b, err := json.Marshal(body)
	if err != nil {
		processingError(w, err, g.ew)
		return
	}
	w.Header().Set(HeaderContentType, HeaderJson)
	w.Write(b)
}
