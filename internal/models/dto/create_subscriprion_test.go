package dto_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Piccadilly98/subscription_service/internal/models/dto"
	"github.com/google/uuid"
)

func TestCreateSubscriptions_Validate(t *testing.T) {
	testCase := []struct {
		name    string
		dto     *dto.CreateSubscriptionsRequest
		wantErr error
	}{
		{
			name: "valid_test_1_not_end_dd-mm-yyyy",
			dto: &dto.CreateSubscriptionsRequest{
				ServiceName: "Yandex Plus",
				Price:       300,
				UserID:      uuid.NewString(),
				StartDate:   "01-01-2025",
			},
		},
		{
			name: "valid_test_2_not_end_mm-yyyy",
			dto: &dto.CreateSubscriptionsRequest{
				ServiceName: "Yandex Plus",
				Price:       300,
				UserID:      uuid.NewString(),
				StartDate:   "03-2025",
			},
		},
		{
			name: "valid_test_3_end_date_dd-mm-yyyy",
			dto: &dto.CreateSubscriptionsRequest{
				ServiceName: "Yandex Plus",
				Price:       300,
				UserID:      uuid.NewString(),
				StartDate:   "01-01-2025",
				EndDate:     getPtrStr("02-02-2025"),
			},
		},
		{
			name: "valid_test_4_end_date_mm-yyyy",
			dto: &dto.CreateSubscriptionsRequest{
				ServiceName: "Yandex Plus",
				Price:       300,
				UserID:      uuid.NewString(),
				StartDate:   "01-2025",
				EndDate:     getPtrStr("02-2025"),
			},
		},

		{
			name: "invalid_1_no_service_name",
			dto: &dto.CreateSubscriptionsRequest{
				Price:     300,
				UserID:    uuid.NewString(),
				StartDate: "01-01-2025",
				EndDate:   getPtrStr("02-02-2025"),
			},
			wantErr: fmt.Errorf("service_name cannot be empty"),
		},
		{
			name: "invalid_2_price<0",
			dto: &dto.CreateSubscriptionsRequest{
				ServiceName: "Yandex Plus",
				Price:       -1,
				UserID:      uuid.NewString(),
				StartDate:   "01-2025",
				EndDate:     getPtrStr("02-2025"),
			},
			wantErr: fmt.Errorf("price cannot be <= 0"),
		},
		{
			name: "invalid_3_price==0",
			dto: &dto.CreateSubscriptionsRequest{
				ServiceName: "Yandex Plus",
				Price:       0,
				UserID:      uuid.NewString(),
				StartDate:   "01-2025",
				EndDate:     getPtrStr("02-2025"),
			},
			wantErr: fmt.Errorf("price cannot be <= 0"),
		},
		{
			name: "invalid_4_id_empty",
			dto: &dto.CreateSubscriptionsRequest{
				ServiceName: "Yandex Plus",
				Price:       300,
				StartDate:   "01-2025",
				EndDate:     getPtrStr("02-2025"),
			},
			wantErr: fmt.Errorf("user_id cannot be empty"),
		},
		{
			name: "invalid_5_user_id_not_uuid",
			dto: &dto.CreateSubscriptionsRequest{
				ServiceName: "Yandex Plus",
				Price:       300,
				UserID:      "asdd",
				StartDate:   "01-2025",
				EndDate:     getPtrStr("02-2025"),
			},
			wantErr: fmt.Errorf("user_id is not uuid"),
		},
		{
			name: "invalid_6_no_start_date",
			dto: &dto.CreateSubscriptionsRequest{
				ServiceName: "Yandex Plus",
				Price:       300,
				UserID:      uuid.NewString(),
			},
			wantErr: fmt.Errorf("start_date cannot be empty"),
		},
		{
			name: "invalid_7_start_date_empty",
			dto: &dto.CreateSubscriptionsRequest{
				ServiceName: "Yandex Plus",
				Price:       300,
				UserID:      uuid.NewString(),
				StartDate:   "",
			},
			wantErr: fmt.Errorf("start_date cannot be empty"),
		},
		{
			name: "invalid_8_end_date_empty",
			dto: &dto.CreateSubscriptionsRequest{
				ServiceName: "Yandex Plus",
				Price:       300,
				UserID:      uuid.NewString(),
				StartDate:   "01-2025",
				EndDate:     getPtrStr(""),
			},
			wantErr: fmt.Errorf("end_date cannot be empty"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.dto.Validate()
			if err != nil {
				if tc.wantErr != nil {
					if tc.wantErr.Error() != err.Error() {
						t.Errorf("ERROR: got: %s, expect: %s\n", err.Error(), tc.wantErr.Error())
					}
				} else {
					t.Errorf("unexpected error: %s\n", err.Error())
				}
			} else {
				if tc.wantErr != nil {
					t.Errorf("ERROR: got: nil, expect: %s\n", tc.wantErr.Error())
				}
			}
		})
	}
}

func TestParseStartDate(t *testing.T) {
	testCases := []struct {
		name       string
		date       string
		wantErr    error
		wantedDate time.Time
	}{
		{
			name:       "valid_1_mm-yyyy",
			date:       "01-2025",
			wantedDate: time.Date(2025, 01, 01, 0, 0, 0, 0, time.UTC),
		},
		{
			name:       "valid_2_dd-mm-yyyy",
			date:       "01-01-2025",
			wantedDate: time.Date(2025, 01, 01, 0, 0, 0, 0, time.UTC),
		},
		{
			name:       "valid_3_dd-mm-yyyy",
			date:       "31-01-2025",
			wantedDate: time.Date(2025, 01, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:       "valid_4_dd-mm-yyyy",
			date:       "28-02-2024",
			wantedDate: time.Date(2024, 02, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name:       "valid_5_dd-mm-yyyy",
			date:       "29-02-2024",
			wantedDate: time.Date(2024, 02, 29, 0, 0, 0, 0, time.UTC),
		},

		{
			name:    "invalid_1_random_str",
			date:    "random",
			wantErr: fmt.Errorf("invalid start_date format"),
		},
		{
			name:    "invalid_2_invalid_date",
			date:    "29-02-2025",
			wantErr: fmt.Errorf("invalid start_date format"),
		},
		{
			name:    "invalid_3_invalid_month",
			date:    "02-14-2025",
			wantErr: fmt.Errorf("invalid start_date format"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			res, err := dto.ParceStartDate(tc.date)
			if err != nil {
				if tc.wantErr != nil {
					if tc.wantErr.Error() != err.Error() {
						t.Errorf("ERROR: got: %s, expect: %s\n", err.Error(), tc.wantErr.Error())
					}
				} else {
					t.Errorf("unexpected error: %s\n", err.Error())
				}
				return
			} else {
				if tc.wantErr != nil {
					t.Errorf("ERROR: got: nil, expect: %s\n", tc.wantErr.Error())
				}
			}

			if res.Compare(tc.wantedDate) != 0 {
				t.Errorf("TIME: got: %s, expect: %s\n", res.String(), tc.wantedDate.String())
			}
		})
	}
}

func TestParceEndDate(t *testing.T) {
	tests := []struct {
		name    string
		date    string
		wantErr bool
		wantDay int
	}{
		{
			name:    "mm-yyyy last_day_month",
			date:    "02-2024",
			wantDay: 29,
		},
		{
			name:    "dd-mm-yyyy",
			date:    "15-02-2024",
			wantDay: 15,
		},
		{
			name:    "december expect 31",
			date:    "12-2024",
			wantDay: 31,
		},
		{
			name:    "january expect 31",
			date:    "01-2024",
			wantDay: 31,
		},
		{
			name:    "april - epxect 30",
			date:    "04-2024",
			wantDay: 30,
		},
		{
			name:    "invalid data",
			date:    "31-02-2024",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dto.ParceEndDate(tt.date)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got == nil {
				t.Fatal("got nil date")
			}

			if got.Day() != tt.wantDay {
				t.Errorf("got day = %d, want %d", got.Day(), tt.wantDay)
			}
		})
	}
}

func TestCreateSubscriptions_ToEntity_Integration(t *testing.T) {
	testCases := []struct {
		name        string
		dto         *dto.CreateSubscriptionsRequest
		wantErr     error
		wantIsEnded bool
		expectedEnd bool
		expectEnd   time.Time
		expectStart time.Time
	}{
		{
			name: "normal_1_is_ended",
			dto: &dto.CreateSubscriptionsRequest{
				ServiceName: "Yandex Plus",
				Price:       300,
				UserID:      uuid.NewString(),
				StartDate:   "01-01-2025",
				EndDate:     getPtrStr("02-01-2025"),
			},
			expectedEnd: true,
			wantIsEnded: true,
			expectEnd:   time.Date(2025, 01, 02, 0, 0, 0, 0, time.UTC),
			expectStart: time.Date(2025, 01, 01, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "normal_2_not_ended",
			dto: &dto.CreateSubscriptionsRequest{
				ServiceName: "WB Club",
				Price:       99,
				UserID:      uuid.NewString(),
				StartDate:   "01-01-2025",
				EndDate:     getPtrStr("02-01-2026"),
			},
			expectedEnd: true,
			wantIsEnded: false,
			expectEnd:   time.Date(2026, 01, 02, 0, 0, 0, 0, time.UTC),
			expectStart: time.Date(2025, 01, 01, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "end_before_start_dd-mm-yyyy",
			dto: &dto.CreateSubscriptionsRequest{
				StartDate: "01-01-2025",
				EndDate:   getPtrStr("31-12-2024"),
			},
			wantErr: fmt.Errorf("end_date cannot be before start_date"),
		},
		{
			name: "end_equal_start_mm-yyyy",
			dto: &dto.CreateSubscriptionsRequest{
				ServiceName: "Yandex Music",
				Price:       399,
				UserID:      uuid.NewString(),
				StartDate:   "01-2025",
				EndDate:     getPtrStr("01-2025"),
			},
			wantIsEnded: true,
			expectedEnd: true,
			expectEnd:   time.Date(2025, 01, 31, 0, 0, 0, 0, time.UTC),
			expectStart: time.Date(2025, 01, 01, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "february_leap_year_mm-yyyy",
			dto: &dto.CreateSubscriptionsRequest{
				ServiceName: "Netflix",
				Price:       599,
				UserID:      uuid.NewString(),
				StartDate:   "02-2024",
				EndDate:     getPtrStr("02-2024"),
			},
			expectedEnd: true,
			wantIsEnded: true,
			expectEnd:   time.Date(2024, 02, 29, 0, 0, 0, 0, time.UTC),
			expectStart: time.Date(2024, 02, 01, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "february_non_leap_year_mm-yyyy",
			dto: &dto.CreateSubscriptionsRequest{
				ServiceName: "Netflix",
				Price:       599,
				UserID:      uuid.NewString(),
				StartDate:   "02-2025",
				EndDate:     getPtrStr("02-2025"),
			},
			expectedEnd: true,
			wantIsEnded: true,
			expectEnd:   time.Date(2025, 02, 28, 0, 0, 0, 0, time.UTC),
			expectStart: time.Date(2025, 02, 01, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "december_to_january_mm-yyyy",
			dto: &dto.CreateSubscriptionsRequest{
				ServiceName: "Spotify",
				Price:       199,
				UserID:      uuid.NewString(),
				StartDate:   "12-2024",
				EndDate:     getPtrStr("01-2026"),
			},
			expectedEnd: true,
			wantIsEnded: false,
			expectEnd:   time.Date(2026, 01, 31, 0, 0, 0, 0, time.UTC),
			expectStart: time.Date(2024, 12, 01, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "year_boundary_dd-mm-yyyy",
			dto: &dto.CreateSubscriptionsRequest{
				ServiceName: "Apple Music",
				Price:       299,
				UserID:      uuid.NewString(),
				StartDate:   "31-12-2024",
				EndDate:     getPtrStr("01-01-2025"),
			},
			expectedEnd: true,
			wantIsEnded: true,
			expectEnd:   time.Date(2025, 01, 01, 0, 0, 0, 0, time.UTC),
			expectStart: time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "mixed_formats_start_dd_mm_yyyy_end_mm_yyyy",
			dto: &dto.CreateSubscriptionsRequest{
				ServiceName: "Yandex Plus",
				Price:       399,
				UserID:      uuid.NewString(),
				StartDate:   "15-07-2024",
				EndDate:     getPtrStr("12-2024"),
			},
			expectedEnd: true,
			wantIsEnded: true,
			expectEnd:   time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
			expectStart: time.Date(2024, 07, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "mixed_formats_start_mm_yyyy_end_dd_mm_yyyy",
			dto: &dto.CreateSubscriptionsRequest{
				ServiceName: "YouTube Premium",
				Price:       349,
				UserID:      uuid.NewString(),
				StartDate:   "06-2024",
				EndDate:     getPtrStr("15-08-2024"),
			},
			expectedEnd: true,
			wantIsEnded: true,
			expectEnd:   time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC),
			expectStart: time.Date(2024, 06, 01, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "no_end_date_lifetime_subscription",
			dto: &dto.CreateSubscriptionsRequest{
				ServiceName: "Microsoft 365",
				Price:       3999,
				UserID:      uuid.NewString(),
				StartDate:   "01-2024",
				EndDate:     nil,
			},
			expectedEnd: false,
			wantIsEnded: false,
			expectStart: time.Date(2024, 01, 01, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "invalid_start_date_yyyy_mm_dd",
			dto: &dto.CreateSubscriptionsRequest{
				ServiceName: "Test",
				Price:       100,
				UserID:      uuid.NewString(),
				StartDate:   "2025-01-01",
			},
			wantErr: fmt.Errorf("invalid start_date format"),
		},
		{
			name: "invalid_end_date_with_slashes",
			dto: &dto.CreateSubscriptionsRequest{
				ServiceName: "Test",
				Price:       100,
				UserID:      uuid.NewString(),
				StartDate:   "01-01-2025",
				EndDate:     getPtrStr("01/01/2025"), // Слеши вместо дефисов
			},
			wantErr: fmt.Errorf("invalid end_date format"),
		},
		{
			name: "end_date_in_wrong_order_yyyy_mm_dd",
			dto: &dto.CreateSubscriptionsRequest{
				ServiceName: "Test",
				Price:       100,
				UserID:      uuid.NewString(),
				StartDate:   "01-01-2025",
				EndDate:     getPtrStr("2025-12-31"), // YYYY-MM-DD
			},
			wantErr: fmt.Errorf("invalid end_date format"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.dto.ToEntitie()
			if err != nil {
				if tc.wantErr != nil {
					if tc.wantErr.Error() != err.Error() {
						t.Errorf("ERROR: got: %s, expect: %s\n", err.Error(), tc.wantErr.Error())
					}
				} else {
					t.Errorf("unexpected error: %s\n", err.Error())
				}
				if res != nil {
					t.Errorf("unepected res\n")
				}
				return
			} else {
				if tc.wantErr != nil {
					t.Errorf("ERROR: got: nil, expect: %s\n", tc.wantErr.Error())
				}
			}
			if res.IsEnded != tc.wantIsEnded {
				t.Errorf("IsEnded: got: %v, expect: %v\n", res.IsEnded, tc.wantIsEnded)
			}
			if res.StartDate.Compare(tc.expectStart) != 0 {
				t.Errorf("Start Date: got: %s, expect: %s\n", res.StartDate.String(), tc.expectStart.String())
			}
			if res.Price != tc.dto.Price {
				t.Errorf("PRICE: got: %d, expect: %d\n", res.Price, tc.dto.Price)
			}
			if res.ServiceName != tc.dto.ServiceName {
				t.Errorf("Service name: got: %s, expect: %s\n", res.ServiceName, tc.dto.ServiceName)
			}
			if res.UserID != tc.dto.UserID {
				t.Errorf("User ID: got: %s, expect: %s\n", res.UserID, tc.dto.UserID)
			}
			if tc.expectedEnd {
				if res.EndDate == nil {
					t.Errorf("END DATE: got: nil, expect: %s\n", tc.expectEnd.String())
				} else {
					if tc.expectEnd.Compare(*res.EndDate) != 0 {
						t.Errorf("End Date: got: %s, expect: %s\n", res.EndDate.String(), tc.expectEnd.String())
					}
				}
			}
		})
	}
}
