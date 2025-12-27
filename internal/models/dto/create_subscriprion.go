package dto

import (
	"fmt"
	"time"

	"github.com/Piccadilly98/subscription_service/internal/models/entities"
	"github.com/google/uuid"
)

const (
	layoutWithDay        = "02-01-2006"
	layoutNotContainsDay = "01-2006"
)

type CreateSubscriptionsRequest struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

func (c *CreateSubscriptionsRequest) Validate() error {
	if c.Price <= 0 {
		return fmt.Errorf("price cannot be <= 0")
	}
	if c.ServiceName == "" {
		return fmt.Errorf("service_name cannot be empty")
	}
	if c.StartDate == "" {
		return fmt.Errorf("start_date cannot be empty")
	}
	if c.UserID == "" {
		return fmt.Errorf("user_id cannot be empty")
	}
	if _, err := uuid.Parse(c.UserID); err != nil {
		return fmt.Errorf("user_id is not uuid")
	}
	if c.EndDate != nil && *c.EndDate == "" {
		return fmt.Errorf("end_date cannot be empty")
	}
	return nil
}

func (c *CreateSubscriptionsRequest) ToEntitie() (*entities.CreateNewSubscriptions, error) {
	var dateStart time.Time
	var endDate *time.Time
	var isEnded bool
	dateStart, err := ParceStartDate(c.StartDate)
	if err != nil {
		return nil, err
	}

	if c.EndDate != nil {
		endDate, err = ParceEndDate(*c.EndDate)
		if err != nil {
			return nil, err
		}

		if endDate.Before(dateStart) {
			return nil, fmt.Errorf("end_date cannot be before start_date")
		}
		if endDate.Before(time.Now()) {
			isEnded = true
		}
	}

	model := &entities.CreateNewSubscriptions{
		UserID:      c.UserID,
		IsEnded:     isEnded,
		Price:       c.Price,
		StartDate:   dateStart,
		EndDate:     endDate,
		ServiceName: c.ServiceName,
	}
	return model, nil
}
