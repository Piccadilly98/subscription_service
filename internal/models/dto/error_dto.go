package dto

import "time"

type ErrorDTO struct {
	Status string `json:"status"`
	Error  string `json:"error_text"`
	Date   string `json:"date"`
}

func NewErrorDto(err error) *ErrorDTO {
	return &ErrorDTO{
		Status: "error",
		Error:  err.Error(),
		Date:   time.Now().String(),
	}
}
