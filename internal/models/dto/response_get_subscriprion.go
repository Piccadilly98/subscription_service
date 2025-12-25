package dto

import (
	"time"

	"github.com/Piccadilly98/subscription_service/internal/models/entities_data_base"
)

const (
	StatusActive     = "active"
	StatusEnded      = "ended"
	StatusNotStarted = "not started"
)

type SubscriptionResponse struct {
	SubscribeID string  `json:"subscribe_id"`
	UserID      string  `json:"user_id"`
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date"`
	Status      string  `json:"status"`
}

func FromEntity(entity *entities_data_base.ReadSubscription) *SubscriptionResponse {
	now := time.Now()
	status := ""
	if entity.StartDate.Before(now) {
		if entity.EndDate != nil && entity.EndDate.Before(now) {
			status = StatusEnded
		} else if entity.EndDate == nil ||
			(entity.EndDate != nil && !entity.EndDate.Before(now)) {
			status = StatusActive
		}
	} else {
		status = StatusNotStarted
	}
	dto := &SubscriptionResponse{
		SubscribeID: entity.SubscribeID,
		UserID:      entity.UserID,
		ServiceName: entity.ServiceName,
		Price:       entity.Price,
		Status:      status,
	}
	return dto
}
