package data_base

import (
	"context"
	"fmt"

	"github.com/Piccadilly98/subscription_service/internal/models/entities"
)

func (d *DataBase) CreateNewSubscribe(ctx context.Context, model *entities.CreateNewSubscriptions) (string, error) {
	id := ""
	err := d.db.QueryRowContext(ctx,
		`INSERT INTO subscriptions(user_id, start_date, end_date, service_name, price, is_ended)
	VALUES($1, $2, $3, $4, $5, $6)
	RETURNING id`,
		model.UserID,
		model.StartDate,
		model.EndDate,
		model.ServiceName,
		model.Price,
		model.IsEnded).Scan(&id)
	return id, err
}

func (d *DataBase) GetSubscriptionByID(ctx context.Context, id string) (*entities.ReadSubscription, error) {
	res := &entities.ReadSubscription{}
	err := d.db.QueryRowContext(ctx,
		`SELECT id, user_id, service_name, price, start_date, end_date FROM subscriptions
		WHERE id = $1;
	`, id).Scan(
		&res.SubscribeID,
		&res.UserID,
		&res.ServiceName,
		&res.Price,
		&res.StartDate,
		&res.EndDate,
	)
	return res, err
}

func (d *DataBase) GetExsistBySubID(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := d.db.QueryRowContext(ctx, `
	SELECT 
	EXISTS(SELECT 1 FROM subscriptions WHERE id = $1)
	FROM subscriptions;`, id).Scan(&exists)

	return exists, err
}

func (d *DataBase) UpdateSubscription(ctx context.Context, entitie *entities.UpdateSubscription, id string) error {
	query, args := d.getQueryAndArgs(entitie, id)

	_, err := d.db.ExecContext(ctx, query, args...)
	return err
}

func (d *DataBase) getQueryAndArgs(entitie *entities.UpdateSubscription, id string) (string, []any) {
	query := "UPDATE subscriptions "
	args := []any{}
	quantityArgs := 0

	if entitie.EndDate != nil {
		quantityArgs++
		query += fmt.Sprintf("SET end_date = $%d", quantityArgs)
		args = append(args, *entitie.EndDate)
	}
	if entitie.Ended != nil {
		quantityArgs++
		if len(args) == 0 {
			query += fmt.Sprintf("SET is_ended = $%d", quantityArgs)
			args = append(args, *entitie.Ended)
		} else {
			query += fmt.Sprintf(",is_ended = $%d", quantityArgs)
			args = append(args, *entitie.Ended)
		}
	}
	if entitie.Price != nil {
		quantityArgs++
		if len(args) == 0 {
			query += fmt.Sprintf("SET price = $%d", quantityArgs)
			args = append(args, *entitie.Price)
		} else {
			query += fmt.Sprintf(",price = $%d", quantityArgs)
			args = append(args, *entitie.Price)
		}
	}
	quantityArgs++
	if len(args) == 0 {
		query += "SET updated_date = NOW()"
	} else {
		query += ", updated_date = NOW()"
	}
	query += fmt.Sprintf(" WHERE id = $%d;", quantityArgs)
	args = append(args, id)

	return query, args
}
