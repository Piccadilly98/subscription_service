package dto

import (
	"fmt"
	"time"

	"github.com/Piccadilly98/subscription_service/internal/models/entities"
)

type UpdateSubscriptionRequest struct {
	Price   *int    `json:"price"`
	EndDate *string `json:"end_date"`
	Ended   *bool   `json:"ended"`
}

func (u *UpdateSubscriptionRequest) Validate() error {
	if u.Price == nil && u.Ended == nil && u.EndDate == nil {
		return fmt.Errorf("not data for update")
	}

	if u.Price != nil {
		if *u.Price <= 0 {
			return fmt.Errorf("price connot be <=0")
		}
	}

	if u.EndDate != nil {
		if *u.EndDate == "" {
			return fmt.Errorf("end_date cannot be empty")
		}
	}

	return nil
}

func (u *UpdateSubscriptionRequest) ToEntitie() (*entities.UpdateSubscription, error) {
	var endDate *time.Time

	if u.EndDate != nil {
		end, err := ParceEndDate(*u.EndDate)
		if err != nil {
			return nil, err
		}
		endDate = end
	}

	return &entities.UpdateSubscription{
		Price:   u.Price,
		EndDate: endDate,
		Ended:   u.Ended,
	}, nil
}
