package dto

import (
	"fmt"

	"github.com/Piccadilly98/subscription_service/internal/models/entities"
)

type SummaryResponse struct {
	Sum       int    `json:"sum" example:"1000"`
	StartDate string `json:"begin_period" example:"08-08-2025"`
	EndDate   string `json:"finish_period" example:"08-08-2026"`
}

func FromEntityToSummaryResponse(entitie *entities.GetSummary, sum int) *SummaryResponse {
	return &SummaryResponse{
		Sum:       sum,
		StartDate: fmt.Sprintf("%02d-%02d-%d", entitie.StartDate.Day(), entitie.StartDate.Month(), entitie.StartDate.Year()),
		EndDate:   fmt.Sprintf("%02d-%02d-%d", entitie.EndDate.Day(), entitie.EndDate.Month(), entitie.EndDate.Year()),
	}
}
