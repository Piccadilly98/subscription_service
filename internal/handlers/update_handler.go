package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Piccadilly98/subscription_service/internal/errors_checker"
	"github.com/Piccadilly98/subscription_service/internal/models/dto"
	"github.com/Piccadilly98/subscription_service/internal/service"
)

type UpdateHandler struct {
	serv *service.Service
	ew   *errors_checker.ErrorWorker
}

func NewUpdateHandler(serv *service.Service, ew *errors_checker.ErrorWorker) *UpdateHandler {
	return &UpdateHandler{
		serv: serv,
		ew:   ew,
	}
}

// Update godoc
// @Summary      Обновление данных подписки
// @Description  Частичное обновление полей подписки по уникальному идентификатору (UUID).<br><br>Сервис выполняет проверки:<br>• ID не пустой и является валидным UUID<br>• Подписка с таким ID существует<br>• Хотя бы одно поле передано для обновления<br>• Изменение завершённой подписки запрещено<br>• Нельзя установить end_date раньше start_date<br>• Нельзя остановить подписку, которая ещё не началась<br>• Цена не может быть меньше или равна нулю<br><br>Опциональные поля в теле запроса:<br>• price — новая месячная стоимость (целое положительное число)<br>• end_date — новая дата окончания (формат "MM-YYYY" или "DD-MM-YYYY", можно null для снятия)<br>• ended — флаг остановки подписки (подписка считается активной в день остановки)<br><br>При успешном обновлении возвращаются актуальные данные подписки.<br><br>Все ошибки возвращаются в едином формате JSON (ErrorDTO).
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id      path  string                    true  "ID подписки (UUID)"  example(b6469ed4-44ae-4436-8827-76131d2d446e)
// @Param        request body dto.UpdateSubscriptionRequest true "Данные для обновления"
// @Success      200 {object} dto.SubscriptionResponse "Данные обновлённой подписки"
// @Failure      400 {object} dto.ErrorDTO "Неверные данные запроса"
// @Failure      404 {object} dto.ErrorDTO "Подписка с указанным ID не найдена"
// @Failure      500 {object} dto.ErrorDTO "Внутренняя ошибка сервера"
// @Router       /subscriptions/{id} [put]
func (u *UpdateHandler) Handler(w http.ResponseWriter, r *http.Request) {
	if !checkHeaderJson(w, r) {
		return
	}
	id := checkURLParam(w, r, u.ew)
	if id == "" {
		return
	}
	exist, err := u.serv.GetExsistBySubID(r.Context(), id)
	if err != nil {
		processingError(w, err, u.ew)
		return
	}
	if !exist {
		processingError(w, errors.New("subscription not found"), u.ew)
		return
	}
	body := &dto.UpdateSubscriptionRequest{}

	err = json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		log.Println(err)
		errorResponse(w, fmt.Errorf("invalid body format"), http.StatusBadRequest)
		return
	}

	res, err := u.serv.UpdateSubscription(r.Context(), body, id)
	if err != nil {
		processingError(w, err, u.ew)
		return
	}

	b, err := json.Marshal(res)
	if err != nil {
		processingError(w, err, u.ew)
		return
	}
	w.Header().Set(HeaderContentType, HeaderJson)
	w.Write(b)
}
