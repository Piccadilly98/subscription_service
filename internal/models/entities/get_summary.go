package entities

import "time"

type GetSummary struct {
	UserID      *string
	ServiceName *string
	StartDate   time.Time
	EndDate     time.Time
}
