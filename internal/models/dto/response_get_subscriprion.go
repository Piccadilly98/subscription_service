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
	SubscribeID string  `json:"subscribe_id" example:"b6469ed4-44ae-4436-8827-76131d2d446e"`
	UserID      string  `json:"user_id"  example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	ServiceName string  `json:"service_name" example:"Yandex Plus"`
	Price       int     `json:"price" example:"400"`
	StartDate   string  `json:"start_date" example:"01-07-2025"`
	EndDate     *string `json:"end_date,omitempty" example:"31-12-2030"`
	Status      string  `json:"status" example:"active"`
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
