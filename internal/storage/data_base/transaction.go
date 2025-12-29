package data_base

import (
	"context"
	"database/sql"

	"github.com/Piccadilly98/subscription_service/internal/models/entities"
)

func (d *DataBase) GetSubscriptionForUpdateTX(ctx context.Context, id string, tx *sql.Tx) (*entities.ReadSubscription, error) {
	res := &entities.ReadSubscription{}
	err := tx.QueryRowContext(ctx,
		`SELECT id, user_id, service_name, price, start_date, end_date, is_ended FROM subscriptions
		WHERE id = $1
		FOR UPDATE;
	`, id).Scan(
		&res.SubscribeID,
		&res.UserID,
		&res.ServiceName,
		&res.Price,
		&res.StartDate,
		&res.EndDate,
		&res.IsEnded,
	)
	return res, err
}

func (d *DataBase) GetSubscriptionTX(ctx context.Context, id string, tx *sql.Tx) (*entities.ReadSubscription, error) {
	res := &entities.ReadSubscription{}
	err := tx.QueryRowContext(ctx,
		`SELECT id, user_id, service_name, price, start_date, end_date, is_ended FROM subscriptions
		WHERE id = $1;
	`, id).Scan(
		&res.SubscribeID,
		&res.UserID,
		&res.ServiceName,
		&res.Price,
		&res.StartDate,
		&res.EndDate,
		&res.IsEnded,
	)
	return res, err
}

func (d *DataBase) UpdateSubscriptionTX(ctx context.Context, entitie *entities.UpdateSubscription, id string, tx *sql.Tx) error {
	query, args := d.getQueryAndArgsUpdate(entitie, id)

	_, err := tx.ExecContext(ctx, query, args...)
	return err
}
