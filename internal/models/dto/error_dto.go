package dto

import "time"

type ErrorDTO struct {
	Status string `json:"status" example:"error"`
	Error  string `json:"error_text"`
	Date   string `json:"date" example:"2025-12-31 00:48:39.814495 +0000 UTC"`
}

func NewErrorDto(err error) *ErrorDTO {
	now := time.Now()
	return &ErrorDTO{
		Status: "error",
		Error:  err.Error(),
		Date:   time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.UTC).String(),
	}
}
