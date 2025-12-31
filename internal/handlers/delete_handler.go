package handlers

import (
	"net/http"

	"github.com/Piccadilly98/subscription_service/internal/errors_checker"
	"github.com/Piccadilly98/subscription_service/internal/service"
)

type DeleteHanlder struct {
	serv *service.Service
	ew   *errors_checker.ErrorWorker
}

func NewDeleteHandler(serv *service.Service, ew *errors_checker.ErrorWorker) *DeleteHanlder {
	return &DeleteHanlder{
		serv: serv,
		ew:   ew,
	}
}

// DeleteSubscription godoc
// @Summary      Удалить существующую подписку
// @Description  Удаляет существующую подписку по уникальному идентификатору (UUID).<br><br>Сервис выполняет проверки:<br>• ID не пустой и является валидным UUID<br>• Подписка с таким ID существует<br><br>При успешном удалении возвращается статус 204 No Content (без тела).<br><br>Все ошибки возвращаются в формате JSON (ErrorDTO).
//
// @Tags         subscriptions
// @Produce      json
// @Param        id  path  string  true  "ID подписки (UUID)"  example(b6469ed4-44ae-4436-8827-76131d2d446e)
// @Success      204 "Подписка успешно удалена"
// @Failure      404 {object} dto.ErrorDTO "ID не существует"
// @Failure      400 {object} dto.ErrorDTO "ID не uuid или пустой"
// @Failure      500 {object} dto.ErrorDTO "Ошибка сервера"
// @Router       /subscriptions/{id} [delete]
func (d *DeleteHanlder) Handler(w http.ResponseWriter, r *http.Request) {
	id := checkURLParam(w, r, d.ew)
	if id == "" {
		return
	}

	err := d.serv.DeleteSubByID(r.Context(), id)
	if err != nil {
		processingError(w, err, d.ew)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
