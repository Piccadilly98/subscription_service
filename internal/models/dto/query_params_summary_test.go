package dto_test

import (
	"fmt"
	"testing"

	"github.com/Piccadilly98/subscription_service/internal/models/dto"
	"github.com/google/uuid"
)

func TestQuerySummary_ToEntitie(t *testing.T) {
	testCases := []struct {
		name          string
		query         *dto.QueryParamsSummary
		expectedError error
	}{
		{
			name: "valid_dd-mm-yyyy",
			query: &dto.QueryParamsSummary{
				UserID:      getPtrStr(uuid.NewString()),
				ServiceName: getPtrStr("random"),
				StartDate:   "01-01-2025",
				EndDate:     "01-01-2026",
			},
		},
		{
			name: "valid_mm-yyyy",
			query: &dto.QueryParamsSummary{
				UserID:      getPtrStr(uuid.NewString()),
				ServiceName: getPtrStr("random"),
				StartDate:   "01-2025",
				EndDate:     "01-2026",
			},
		},
		{
			name: "valid_no_service_name",
			query: &dto.QueryParamsSummary{
				UserID:    getPtrStr(uuid.NewString()),
				StartDate: "01-01-2025",
				EndDate:   "01-01-2026",
			},
		},
		{
			name: "valid_no_user_id",
			query: &dto.QueryParamsSummary{
				ServiceName: getPtrStr("random"),
				StartDate:   "01-01-2025",
				EndDate:     "01-01-2026",
			},
		},
		{
			name: "valid_no_user_id_and_service_name",
			query: &dto.QueryParamsSummary{
				StartDate: "01-01-2025",
				EndDate:   "01-01-2026",
			},
		},
		{
			name: "invalid_user_id_not_uuid",
			query: &dto.QueryParamsSummary{
				UserID:      getPtrStr("random"),
				ServiceName: getPtrStr("random"),
				StartDate:   "01-01-2025",
				EndDate:     "01-01-2026",
			},
			expectedError: fmt.Errorf("[USER]user_id is not uuid"),
		},
		{
			name: "invalid_user_id_empty",
			query: &dto.QueryParamsSummary{
				UserID:      getPtrStr(""),
				ServiceName: getPtrStr("random"),
				StartDate:   "01-01-2025",
				EndDate:     "01-01-2026",
			},
			expectedError: fmt.Errorf("[USER]user_id is not uuid"),
		},
		{
			name: "invalid_no_end_date",
			query: &dto.QueryParamsSummary{
				UserID:      getPtrStr(uuid.NewString()),
				ServiceName: getPtrStr("random"),
				StartDate:   "01-01-2025",
			},
			expectedError: fmt.Errorf("[USER]no contains end_period_date"),
		},
		{
			name: "invalid_no_start_date",
			query: &dto.QueryParamsSummary{
				UserID:      getPtrStr(uuid.NewString()),
				ServiceName: getPtrStr("random"),
				StartDate:   "",
				EndDate:     "01-01-2025",
			},
			expectedError: fmt.Errorf("[USER]no contains start_period_date"),
		},
		{
			name: "invalid_no_dates",
			query: &dto.QueryParamsSummary{
				UserID:      getPtrStr(uuid.NewString()),
				ServiceName: getPtrStr("random"),
			},
			expectedError: fmt.Errorf("[USER]no contains period dates"),
		},
		{
			name: "invalid_invalid_format_start_date_dd-mm-yyyy",
			query: &dto.QueryParamsSummary{
				UserID:      getPtrStr(uuid.NewString()),
				ServiceName: getPtrStr("random"),
				StartDate:   "01/01/2000",
				EndDate:     "01-01-2026",
			},
			expectedError: fmt.Errorf("[USER]invalid start_date format"),
		},
		{
			name: "invalid_invalid_format_start_date_mm-yyyy",
			query: &dto.QueryParamsSummary{
				UserID:      getPtrStr(uuid.NewString()),
				ServiceName: getPtrStr("random"),
				StartDate:   "01/2000",
				EndDate:     "01-01-2026",
			},
			expectedError: fmt.Errorf("[USER]invalid start_date format"),
		},
		{
			name: "invalid_invalid_format_start_date_random",
			query: &dto.QueryParamsSummary{
				UserID:      getPtrStr(uuid.NewString()),
				ServiceName: getPtrStr("random"),
				StartDate:   "random",
				EndDate:     "01-01-2026",
			},
			expectedError: fmt.Errorf("[USER]invalid start_date format"),
		},

		{
			name: "invalid_invalid_format_end_date_dd-mm-yyyy",
			query: &dto.QueryParamsSummary{
				UserID:      getPtrStr(uuid.NewString()),
				ServiceName: getPtrStr("random"),
				StartDate:   "01-01-2000",
				EndDate:     "01/01/2026",
			},
			expectedError: fmt.Errorf("[USER]invalid end_date format"),
		},
		{
			name: "invalid_invalid_format_end_date_mm-yyyy",
			query: &dto.QueryParamsSummary{
				UserID:      getPtrStr(uuid.NewString()),
				ServiceName: getPtrStr("random"),
				StartDate:   "01-01-2000",
				EndDate:     "01/2026",
			},
			expectedError: fmt.Errorf("[USER]invalid end_date format"),
		},
		{
			name: "invalid_invalid_format_end_date_random",
			query: &dto.QueryParamsSummary{
				UserID:      getPtrStr(uuid.NewString()),
				ServiceName: getPtrStr("random"),
				StartDate:   "14-12-2025",
				EndDate:     "random",
			},
			expectedError: fmt.Errorf("[USER]invalid end_date format"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.query.ToEntitie()
			if err != nil {
				if tc.expectedError != nil {
					if tc.expectedError.Error() != err.Error() {
						t.Errorf("ERROR: got: %s, expect: %s\n", err.Error(), tc.expectedError.Error())
					}
				} else {
					t.Errorf("unexpected error: %s\n", err.Error())
				}
				return
			} else {
				if tc.expectedError != nil {
					t.Errorf("ERROR: got: nil, expect: %snil\n", tc.expectedError.Error())
				}
			}

			if res.ServiceName != nil {
				if tc.query.ServiceName != nil {
					if *tc.query.ServiceName != *res.ServiceName {
						t.Errorf("SERVICE NAME: got: %s, expect: %s\n", *res.ServiceName, *tc.query.ServiceName)
					}
				} else {
					t.Errorf("unexpected service name: %s\n", *res.ServiceName)
				}
			}

			if res.UserID != nil {
				if tc.query.UserID != nil {
					if *tc.query.UserID != *res.UserID {
						t.Errorf("USER_ID: got: %s, expect: %s\n", *res.UserID, *tc.query.UserID)
					}
				} else {
					t.Errorf("unexpected user_id: %s\n", *res.UserID)
				}
			}
		})
	}
}
