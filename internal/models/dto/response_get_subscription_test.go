package dto_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Piccadilly98/subscription_service/internal/models/dto"
	"github.com/Piccadilly98/subscription_service/internal/models/entities"
	"github.com/google/uuid"
)

func TestSubscriptionResponse_FromEntityToSubResp(t *testing.T) {
	Now := time.Now()
	testCase := []struct {
		name            string
		entitie         *entities.ReadSubscription
		expectStatus    string
		expectedEndDate bool
		expectEndDate   string
		expectStartDate string
	}{
		{
			name: "end_date_nil_status_not_started",
			entitie: &entities.ReadSubscription{
				SubscribeID: uuid.NewString(),
				UserID:      uuid.NewString(),
				ServiceName: "Yandex Plus",
				Price:       400,
				StartDate:   Now.Add(24 * time.Hour),
				EndDate:     nil,
				IsEnded:     false,
			},
			expectStatus:    dto.StatusNotStarted,
			expectedEndDate: false,
			expectStartDate: fmt.Sprintf("%02d-%02d-%02d",
				Now.Add(24*time.Hour).Day(),
				Now.Add(24*time.Hour).Month(),
				Now.Add(24*time.Hour).Year()),
		},
		{
			name: "end_date_nil_status_active",
			entitie: &entities.ReadSubscription{
				SubscribeID: uuid.NewString(),
				UserID:      uuid.NewString(),
				ServiceName: "Yandex Plus",
				Price:       400,
				StartDate:   Now.Add(-24 * time.Hour),
				EndDate:     nil,
				IsEnded:     false,
			},
			expectStatus:    dto.StatusActive,
			expectedEndDate: false,
			expectStartDate: fmt.Sprintf("%02d-%02d-%02d",
				Now.Add(-24*time.Hour).Day(),
				Now.Add(-24*time.Hour).Month(),
				Now.Add(-24*time.Hour).Year()),
		},
		{
			name: "end_date_not_nil_status_ended",
			entitie: &entities.ReadSubscription{
				SubscribeID: uuid.NewString(),
				UserID:      uuid.NewString(),
				ServiceName: "Yandex Plus",
				Price:       400,
				StartDate:   Now.Add(-72 * time.Hour),
				EndDate: GetTimeDate(time.Date(Now.Add(-2*time.Hour).Year(),
					Now.Add(-24*time.Hour).Month(),
					Now.Add(-24*time.Hour).Day(),
					0, 0, 0, 0, time.UTC)),
				IsEnded: false,
			},
			expectStatus:    dto.StatusEnded,
			expectedEndDate: true,
			expectEndDate: fmt.Sprintf("%02d-%02d-%02d",
				Now.Add(-24*time.Hour).Day(),
				Now.Add(-24*time.Hour).Month(),
				Now.Add(-24*time.Hour).Year()),
			expectStartDate: fmt.Sprintf("%02d-%02d-%02d",
				Now.Add(-72*time.Hour).Day(),
				Now.Add(-72*time.Hour).Month(),
				Now.Add(-72*time.Hour).Year()),
		},
		{
			name: "end_date_not_nil_status_not_started",
			entitie: &entities.ReadSubscription{
				SubscribeID: uuid.NewString(),
				UserID:      uuid.NewString(),
				ServiceName: "Yandex Plus",
				Price:       400,
				StartDate:   Now.Add(24 * time.Hour),
				EndDate: GetTimeDate(time.Date(
					Now.Add(48*time.Hour).Year(),
					Now.Add(48*time.Hour).Month(),
					Now.Add(48*time.Hour).Day(),
					0, 0, 0, 0, time.UTC)),
				IsEnded: false,
			},
			expectStatus:    dto.StatusNotStarted,
			expectedEndDate: true,
			expectEndDate: fmt.Sprintf("%02d-%02d-%d",
				Now.Add(48*time.Hour).Day(),
				Now.Add(48*time.Hour).Month(),
				Now.Add(48*time.Hour).Year()),
			expectStartDate: fmt.Sprintf("%02d-%02d-%02d",
				Now.Add(24*time.Hour).Day(),
				Now.Add(24*time.Hour).Month(),
				Now.Add(24*time.Hour).Year()),
		},
		{
			name: "end_date_not_nil_status_ended",
			entitie: &entities.ReadSubscription{
				SubscribeID: uuid.NewString(),
				UserID:      uuid.NewString(),
				ServiceName: "Yandex Plus",
				Price:       400,
				StartDate:   Now.Add(-48 * time.Hour),
				EndDate: GetTimeDate(time.Date(
					Now.Add(-24*time.Hour).Year(),
					Now.Add(-24*time.Hour).Month(),
					Now.Add(-24*time.Hour).Day(),
					0, 0, 0, 0, time.UTC)),
				IsEnded: false,
			},
			expectStatus:    dto.StatusEnded,
			expectedEndDate: true,
			expectEndDate: fmt.Sprintf("%02d-%02d-%02d",
				Now.Add(-24*time.Hour).Day(),
				Now.Add(-24*time.Hour).Month(),
				Now.Add(-24*time.Hour).Year()),
			expectStartDate: fmt.Sprintf("%02d-%02d-%02d",
				Now.Add(-48*time.Hour).Day(),
				Now.Add(-48*time.Hour).Month(),
				Now.Add(-48*time.Hour).Year()),
		},
		{
			name: "starts_today_no_end_date",
			entitie: &entities.ReadSubscription{
				SubscribeID: uuid.NewString(),
				UserID:      uuid.NewString(),
				ServiceName: "Service",
				Price:       100,
				StartDate:   Now,
				EndDate:     nil,
				IsEnded:     false,
			},
			expectStatus:    dto.StatusActive,
			expectedEndDate: false,
			expectStartDate: fmt.Sprintf("%02d-%02d-%d",
				Now.Day(),
				Now.Month(),
				Now.Year()),
		},
		{
			name: "ended_today",
			entitie: &entities.ReadSubscription{
				SubscribeID: uuid.NewString(),
				UserID:      uuid.NewString(),
				ServiceName: "Service",
				Price:       100,
				StartDate:   Now.Add(-48 * time.Hour),
				EndDate: GetTimeDate(time.Date(
					Now.Year(),
					Now.Month(),
					Now.Day(), 0, 0, 0, 0, time.UTC)),
				IsEnded: false,
			},
			expectStatus:    dto.StatusActive,
			expectedEndDate: true,
			expectEndDate:   fmt.Sprintf("%02d-%02d-%d", Now.Day(), Now.Month(), Now.Year()),
			expectStartDate: fmt.Sprintf("%02d-%02d-%d",
				Now.Add(-48*time.Hour).Day(),
				Now.Add(-48*time.Hour).Month(),
				Now.Add(-48*time.Hour).Year()),
		},
		{
			name: "ends_last_day_of_month",
			entitie: &entities.ReadSubscription{
				SubscribeID: uuid.NewString(),
				UserID:      uuid.NewString(),
				ServiceName: "Service",
				Price:       100,
				StartDate:   Now.Add(-24 * time.Hour),
				EndDate: GetTimeDate(time.Date(
					Now.Year(), Now.Month()+1, 0, 0, 0, 0, 0, time.UTC)),
				IsEnded: false,
			},
			expectStatus:    dto.StatusActive,
			expectedEndDate: true,
			expectEndDate: fmt.Sprintf("%02d-%02d-%d",
				time.Date(Now.Year(), Now.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day(),
				Now.Month(),
				Now.Year()),
			expectStartDate: fmt.Sprintf("%02d-%02d-%d",
				Now.Add(-24*time.Hour).Day(),
				Now.Add(-24*time.Hour).Month(),
				Now.Add(-24*time.Hour).Year()),
		},
		{
			name: "is_ended_true_but_end_date_future",
			entitie: &entities.ReadSubscription{
				SubscribeID: uuid.NewString(),
				UserID:      uuid.NewString(),
				ServiceName: "Service",
				Price:       100,
				StartDate:   Now.Add(-24 * time.Hour),
				EndDate:     GetTimeDate(Now.Add(24 * time.Hour)),
				IsEnded:     true,
			},
			expectStatus:    dto.StatusActive,
			expectedEndDate: true,
			expectEndDate: fmt.Sprintf("%02d-%02d-%d",
				Now.Add(24*time.Hour).Day(),
				Now.Add(24*time.Hour).Month(),
				Now.Add(24*time.Hour).Year()),
			expectStartDate: fmt.Sprintf("%02d-%02d-%d",
				Now.Add(-24*time.Hour).Day(),
				Now.Add(-24*time.Hour).Month(),
				Now.Add(-24*time.Hour).Year()),
		},
		{
			name: "same_start_and_end_date",
			entitie: &entities.ReadSubscription{
				SubscribeID: uuid.NewString(),
				UserID:      uuid.NewString(),
				ServiceName: "One Day Service",
				Price:       50,
				StartDate:   Now.Add(-24 * time.Hour),
				EndDate:     GetTimeDate(Now.Add(-24 * time.Hour)),
				IsEnded:     false,
			},
			expectStatus:    dto.StatusEnded,
			expectedEndDate: true,
			expectEndDate: fmt.Sprintf("%02d-%02d-%d",
				Now.Add(-24*time.Hour).Day(),
				Now.Add(-24*time.Hour).Month(),
				Now.Add(-24*time.Hour).Year()),
			expectStartDate: fmt.Sprintf("%02d-%02d-%d",
				Now.Add(-24*time.Hour).Day(),
				Now.Add(-24*time.Hour).Month(),
				Now.Add(-24*time.Hour).Year()),
		},
		{
			name: "not_started_with_end_date",
			entitie: &entities.ReadSubscription{
				SubscribeID: uuid.NewString(),
				UserID:      uuid.NewString(),
				ServiceName: "Future Service",
				Price:       150,
				StartDate:   Now.Add(72 * time.Hour),
				EndDate:     GetTimeDate(Now.Add(96 * time.Hour)),
				IsEnded:     false,
			},
			expectStatus:    dto.StatusNotStarted,
			expectedEndDate: true,
			expectEndDate: fmt.Sprintf("%02d-%02d-%d",
				Now.Add(96*time.Hour).Day(),
				Now.Add(96*time.Hour).Month(),
				Now.Add(96*time.Hour).Year()),
			expectStartDate: fmt.Sprintf("%02d-%02d-%d",
				Now.Add(72*time.Hour).Day(),
				Now.Add(72*time.Hour).Month(),
				Now.Add(72*time.Hour).Year()),
		},
		{
			name: "end_date_before_start_date",
			entitie: &entities.ReadSubscription{
				SubscribeID: uuid.NewString(),
				UserID:      uuid.NewString(),
				ServiceName: "Invalid Service",
				Price:       100,
				StartDate:   Now.Add(24 * time.Hour),
				EndDate:     GetTimeDate(Now.Add(-24 * time.Hour)),
				IsEnded:     false,
			},
			expectStatus:    dto.StatusNotStarted,
			expectedEndDate: true,
			expectEndDate: fmt.Sprintf("%02d-%02d-%d",
				Now.Add(-24*time.Hour).Day(),
				Now.Add(-24*time.Hour).Month(),
				Now.Add(-24*time.Hour).Year()),
			expectStartDate: fmt.Sprintf("%02d-%02d-%d",
				Now.Add(24*time.Hour).Day(),
				Now.Add(24*time.Hour).Month(),
				Now.Add(24*time.Hour).Year()),
		},
		{
			name: "zero_price",
			entitie: &entities.ReadSubscription{
				SubscribeID: uuid.NewString(),
				UserID:      uuid.NewString(),
				ServiceName: "Free Service",
				Price:       0,
				StartDate:   Now.Add(-24 * time.Hour),
				EndDate:     nil,
				IsEnded:     false,
			},
			expectStatus:    dto.StatusActive,
			expectedEndDate: false,
			expectStartDate: fmt.Sprintf("%02d-%02d-%d",
				Now.Add(-24*time.Hour).Day(),
				Now.Add(-24*time.Hour).Month(),
				Now.Add(-24*time.Hour).Year()),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			resp := dto.FromEntityToSubResp(tc.entitie)
			if resp.SubscribeID != tc.entitie.SubscribeID {
				t.Errorf("SUBSCRIBE ID: got: %s, expect: %s\n", resp.SubscribeID, tc.entitie.SubscribeID)
			}
			if resp.UserID != tc.entitie.UserID {
				t.Errorf("USER ID: got: %s, expect: %s\n", resp.UserID, tc.entitie.UserID)
			}
			if resp.ServiceName != tc.entitie.ServiceName {
				t.Errorf("SERVICE NAME: got: %s, expect: %s\n", resp.ServiceName, tc.entitie.ServiceName)
			}
			if resp.Price != tc.entitie.Price {
				t.Errorf("PRICE: got: %d, expect: %d\n", resp.Price, tc.entitie.Price)
			}
			if resp.StartDate != tc.expectStartDate {
				t.Errorf("START DATE: got: %s, expect: %s\n", resp.StartDate, tc.expectStartDate)
			}
			if tc.expectStatus != resp.Status {
				t.Errorf("STATUS: got: %s, expect: %s\n", resp.Status, tc.expectStatus)
			}
			if tc.expectedEndDate {
				if resp.EndDate != nil {
					if tc.expectEndDate != *resp.EndDate {
						t.Errorf("END DATE: got: %s, expect: %s\n", *resp.EndDate, tc.expectEndDate)
					}
				} else {
					t.Errorf("END DATE: got: nil, expect: %s\n", tc.expectEndDate)
				}
			} else {
				if resp.EndDate != nil {
					t.Errorf("unexpected endDate: %s\n", *resp.EndDate)
				}
			}
		})
	}
}

func GetTimeDate(t time.Time) *time.Time {
	return &t
}
