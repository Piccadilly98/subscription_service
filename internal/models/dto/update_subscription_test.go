package dto_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Piccadilly98/subscription_service/internal/models/dto"
)

func TestUpdateSubscription_Validate(t *testing.T) {
	testCase := []struct {
		name          string
		dto           *dto.UpdateSubscriptionRequest
		ExpectedError error
	}{
		{
			name: "valid_end_date_ended_nil",
			dto: &dto.UpdateSubscriptionRequest{
				EndDate: getPtrStr("02-02-2003"),
			},
		},
		{
			name: "valid_end_date_ended_nil_price_not_nil",
			dto: &dto.UpdateSubscriptionRequest{
				EndDate: getPtrStr("02-02-2003"),
				Price:   getIntPtr(300),
			},
		},
		{
			name: "end_not_nil_and_ended_true",
			dto: &dto.UpdateSubscriptionRequest{
				EndDate: getPtrStr("02-02-2003"),
				Ended:   getBoolPtr(true),
			},
			ExpectedError: fmt.Errorf("[USER]cannot specify both end_date and ended"),
		},
		{
			name: "end_not_nil_and_ended_false",
			dto: &dto.UpdateSubscriptionRequest{
				EndDate: getPtrStr("02-02-2003"),
				Ended:   getBoolPtr(false),
			},
			ExpectedError: fmt.Errorf("[USER]cannot specify both end_date and ended"),
		},
		{
			name: "price<0",
			dto: &dto.UpdateSubscriptionRequest{
				Price: getIntPtr(-1),
			},
			ExpectedError: fmt.Errorf("[USER]price connot be <=0"),
		},
		{
			name: "price==0",
			dto: &dto.UpdateSubscriptionRequest{
				Price: getIntPtr(-1),
			},
			ExpectedError: fmt.Errorf("[USER]price connot be <=0"),
		},
		{
			name:          "not_update_data",
			dto:           &dto.UpdateSubscriptionRequest{},
			ExpectedError: fmt.Errorf("[USER]not data for update"),
		},
		{
			name: "end_date_empty",
			dto: &dto.UpdateSubscriptionRequest{
				EndDate: getPtrStr(""),
			},
			ExpectedError: fmt.Errorf("[USER]end_date cannot be empty"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.dto.Validate()
			if err != nil {
				if tc.ExpectedError != nil {
					if tc.ExpectedError.Error() != err.Error() {
						t.Errorf("ERROR: got: %s, expect: %s\n", err.Error(), tc.ExpectedError.Error())
					}
				} else {
					t.Errorf("unexpected error: %s\n", err.Error())
				}
			} else {
				if tc.ExpectedError != nil {
					t.Errorf("ERROR: got: nil, expect: %s\n", tc.ExpectedError.Error())
				}
			}
		})
	}
}

func TestUpdateSubscription_ToEntitie(t *testing.T) {
	testCase := []struct {
		name            string
		dto             *dto.UpdateSubscriptionRequest
		ExpectedEndDate bool
		ExpectEndDate   time.Time
		ExpectError     error
	}{
		{
			name: "valid_test_end_date_mm-yyyy",
			dto: &dto.UpdateSubscriptionRequest{
				EndDate: getPtrStr("01-2025"),
			},
			ExpectedEndDate: true,
			ExpectEndDate:   *GetTimeDate(time.Date(2025, 01, 31, 0, 0, 0, 0, time.UTC)),
		},
		{
			name: "valid_test_end_date_dd-mm-yyyy",
			dto: &dto.UpdateSubscriptionRequest{
				EndDate: getPtrStr("01-01-2025"),
			},
			ExpectedEndDate: true,
			ExpectEndDate:   *GetTimeDate(time.Date(2025, 01, 01, 0, 0, 0, 0, time.UTC)),
		},
		{
			name: "random_end_date",
			dto: &dto.UpdateSubscriptionRequest{
				EndDate: getPtrStr("random"),
			},
			ExpectedEndDate: false,
			ExpectError:     fmt.Errorf("[USER]invalid end_date format"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.dto.ToEntitie()
			if err != nil {
				if tc.ExpectError != nil {
					if tc.ExpectError.Error() != err.Error() {
						t.Errorf("ERROR: got: %s, expect: %s\n", err.Error(), tc.ExpectError.Error())
					}
				} else {
					t.Errorf("unexpected error: %s\n", err.Error())
				}
				return
			} else {
				if tc.ExpectError != nil {
					t.Errorf("ERROR: got: nil, expect: %s\n", tc.ExpectError.Error())
				}
			}
			if res.Ended != nil {
				if tc.dto.Ended != nil {
					if *res.Ended != *tc.dto.Ended {
						t.Errorf("ENDED: got %v, expect: %v\n", *res.Ended, *tc.dto.Ended)
					}
				} else {
					t.Errorf("ENDED: unexpected ended: %v\n", *res.Ended)
				}
			}
			if res.Price != nil {
				if tc.dto.Price != nil {
					if *res.Price != *tc.dto.Price {
						t.Errorf("PRICE: got %d, expect: %d\n", *res.Price, *tc.dto.Price)
					}
				} else {
					t.Errorf("PRICE: unexpected price: %d\n", *res.Price)
				}
			}
			if res.EndDate != nil {
				if tc.ExpectedEndDate {
					if tc.ExpectEndDate.Compare(*res.EndDate) != 0 {
						t.Errorf("END_DATE: got: %s, expect: %s\n", res.EndDate.String(), tc.ExpectEndDate.String())
					}
				} else {
					t.Errorf("END_DATE: unexpectde end_date: %s\n", res.EndDate.String())
				}
			}
		})
	}
}
