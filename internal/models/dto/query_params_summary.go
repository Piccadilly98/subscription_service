package dto

import (
	"fmt"

	"github.com/Piccadilly98/subscription_service/internal/models/entities"
	"github.com/google/uuid"
)

type QueryParamsSummary struct {
	UserID      *string
	ServiceName *string
	StartDate   string
	EndDate     string
}

func (q *QueryParamsSummary) ToEntitie() (*entities.GetSummary, error) {
	res := &entities.GetSummary{}
	if q.EndDate == "" && q.StartDate == "" {
		return nil, fmt.Errorf("[USER]no contains period dates")
	}
	if q.StartDate == "" {
		return nil, fmt.Errorf("[USER]no contains start_period_date")
	}
	if q.EndDate == "" {
		return nil, fmt.Errorf("[USER]no contains end_period_date")
	}
	start, err := ParceStartDate(q.StartDate)
	if err != nil {
		return nil, err
	}
	res.StartDate = start

	end, err := ParceEndDate(q.EndDate)
	if err != nil {
		return nil, err
	}
	res.EndDate = *end
	if q.UserID != nil {
		_, err := uuid.Parse(*q.UserID)
		if err != nil {
			return nil, fmt.Errorf("[USER]user_id is not uuid")
		}
		res.UserID = q.UserID
	}
	res.ServiceName = q.ServiceName

	return res, nil
}
