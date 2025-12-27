package dto

import (
	"fmt"
	"time"

	"github.com/Piccadilly98/subscription_service/internal/models/entities"
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

func FromEntityToSubResp(entity *entities.ReadSubscription) *SubscriptionResponse {
	status := processingStatus(entity)

	var dateEnd *string
	dateStart := ""
	if entity.EndDate != nil {
		end := fmt.Sprintf("%02d-%02d-%d", entity.EndDate.Day(), entity.EndDate.Month(), entity.EndDate.Year())
		dateEnd = &end
	}
	dateStart = fmt.Sprintf("%02d-%02d-%d", entity.StartDate.Day(), entity.StartDate.Month(), entity.StartDate.Year())
	dto := &SubscriptionResponse{
		SubscribeID: entity.SubscribeID,
		UserID:      entity.UserID,
		ServiceName: entity.ServiceName,
		Price:       entity.Price,
		Status:      status,
		EndDate:     dateEnd,
		StartDate:   dateStart,
	}
	return dto
}

func processingStatus(entity *entities.ReadSubscription) string {
	nowTime := time.Now()
	now := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 0, 0, 0, 0, time.UTC)
	startDateDay := time.Date(entity.StartDate.Year(), entity.StartDate.Month(), entity.StartDate.Day(), 0, 0, 0, 0, time.UTC)
	if startDateDay.After(now) {
		return StatusNotStarted
	}

	if entity.EndDate != nil {
		if entity.EndDate.Before(now) {
			return StatusEnded
		} else {
			return StatusActive
		}
	} else {
		return StatusActive
	}
}
