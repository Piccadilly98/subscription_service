package entities_data_base

import "time"

type CreateNewSubscriptions struct {
	UserID      string
	ServiceName string
	Price       int
	StartDate   time.Time
	EndDate     *time.Time
	IsEnded     bool
}
