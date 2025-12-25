package dto

import (
	"fmt"
	"time"

	"github.com/Piccadilly98/subscription_service/internal/models/entities_data_base"
	"github.com/google/uuid"
)

const (
	layoutWithDay        = "02-01-2006"
	layoutNotContainsDay = "01-2006"
)

type RequestCreateSubscriptions struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

func (c *RequestCreateSubscriptions) Validate() error {
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

func (c *RequestCreateSubscriptions) ToEntitie() (*entities_data_base.CreateNewSubscriptions, error) {
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
		if endDate.Compare(dateStart) == 0 {
			return nil, fmt.Errorf("end_date most be == start_date")
		}
		if endDate.Before(time.Now()) {
			isEnded = true
		}
	}

	model := &entities_data_base.CreateNewSubscriptions{
		UserID:      c.UserID,
		IsEnded:     isEnded,
		Price:       c.Price,
		StartDate:   dateStart,
		EndDate:     endDate,
		ServiceName: c.ServiceName,
	}
	return model, nil
}

func ParceStartDate(date string) (time.Time, error) {
	dateStart, err := time.Parse(layoutWithDay, date)
	if err != nil {
		dateStart, err = time.Parse(layoutNotContainsDay, date)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid start_date format")
		}
		date := time.Date(dateStart.Year(), dateStart.Month(), 1, 0, 0, 0, 0, dateStart.Location())
		dateStart = date
	}
	return dateStart, nil
}

func ParceEndDate(date string) (*time.Time, error) {
	endDate, err := time.Parse(layoutWithDay, date)
	if err != nil {
		endDate, err = time.Parse(layoutNotContainsDay, date)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date format")
		}
		end := time.Date(endDate.Year(), endDate.Month()+1, endDate.Day()-1, 0, 0, 0, 0, endDate.Location())
		endDate = end
	}
	return &endDate, nil
}
