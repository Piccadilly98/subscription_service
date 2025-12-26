package entities

import "time"

type UpdateSubscription struct {
	Price   *int
	EndDate *time.Time
	Ended   *bool
}
