package data_base

import (
	"context"

	"github.com/Piccadilly98/subscription_service/internal/models/entities_data_base"
)

func (d *DataBase) CreateNewSubscribe(ctx context.Context, model *entities_data_base.CreateNewSubscriptions) error {
	_, err := d.db.ExecContext(ctx,
		`INSERT INTO subscriptions(user_id, start_date, end_date, service_name, price, is_ended)
	VALUES($1, $2, $3, $4, $5, $6);`,
		model.UserID,
		model.StartDate,
		model.EndDate,
		model.ServiceName,
		model.Price,
		model.IsEnded)
	return err
}
