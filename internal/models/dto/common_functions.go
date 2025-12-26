package dto

import (
	"fmt"
	"time"
)

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
