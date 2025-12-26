package service

import "time"

func GetTimePtr(t time.Time) *time.Time {
	return &t
}

func getIntPtr(i int) *int {
	return &i
}

func getBoolPtr(b bool) *bool {
	return &b
}

func getPtrStr(str string) *string {
	return &str
}
