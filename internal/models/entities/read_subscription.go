package entities

import (
	"time"
)

type ReadSubscription struct {
	SubscribeID string
	UserID      string
	ServiceName string
	Price       int
	StartDate   time.Time
	EndDate     *time.Time
	IsEnded     bool
}
