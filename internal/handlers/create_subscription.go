package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Piccadilly98/subscription_service/internal/errors_checker"
	dto "github.com/Piccadilly98/subscription_service/internal/models/dto"
	"github.com/Piccadilly98/subscription_service/internal/service"
)

type CreateSubsriptionHandler struct {
	service *service.Service
	ew      *errors_checker.ErrorWorker
}

func NewCreateHandler(s *service.Service, ew *errors_checker.ErrorWorker) *CreateSubsriptionHandler {
	return &CreateSubsriptionHandler{
		service: s,
		ew:      ew,
	}
}

// CreateSubscription godoc
// @Summary      Создать новую подписку и получить созданную подписку
// @Description  Принимает данные о новой подписке пользователя в формате JSON и выполняет полную валидацию.<br><br>Обязательные поля: service_name, price, user_id, start_date.<br>Поле end_date — опциональное (можно передать null или не указывать).<br><br>Поддерживаемые форматы даты:<br>• "MM-YYYY" (например, "07-2025")<br>• "DD-MM-YYYY" (например, "15-07-2025")<br><br>Хендлер автоматически:<br>• валидирует корректность дат<br>• проверяет, что end_date не раньше start_date (если указана)<br>• приводит обе даты к единому формату "DD-MM-YYYY" (первый день месяца(для стартовой даты и последний день месяца для даты окончания) для формата "MM-YYYY")<br><br>При успешном создании возвращает данные подписки с присвоенным ID и статусом "active", "ended" или "not started" в зависимости от дат которые были в теле запроса<br><br>При ошибке возвращается dto с текстом ошибки и временем ответа.
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateSubscriptionsRequest true "Данные подписки"
// @Success      201 {object} dto.SubscriptionResponse "Подписка успешно создана"
// @Failure      400 {object} dto.ErrorDTO "Неверные данные запроса"
// @Failure      500 {object} dto.ErrorDTO "Внутренняя ошибка сервера"
// @Router       /subscriptions [post]
func (c *CreateSubsriptionHandler) Handler(w http.ResponseWriter, r *http.Request) {
	if !checkHeaderJson(w, r) {
		return
	}

	body := &dto.CreateSubscriptionsRequest{}
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		log.Println(err)
		errorResponse(w, fmt.Errorf("invalid body format"), http.StatusBadRequest)
		return
	}

	res, err := c.service.CreateSubsription(r.Context(), body)
	if err != nil {
		processingError(w, err, c.ew)
		return
	}

	b, err := json.Marshal(res)
	if err != nil {
		processingError(w, err, c.ew)
		return
	}

	w.Header().Set(HeaderContentType, HeaderJson)
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}
