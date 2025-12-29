package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Piccadilly98/subscription_service/internal/models/dto"
)

func TestErrorResponse(t *testing.T) {
	testCases := []struct {
		name                   string
		err                    error
		code                   int
		expectCode             int
		expectHederContentType string
		expectErrorStr         string
	}{
		{
			name:                   "invalid_syntax",
			err:                    fmt.Errorf("invalid_syntax"),
			code:                   http.StatusBadRequest,
			expectErrorStr:         "invalid_syntax",
			expectHederContentType: HeaderJson,
			expectCode:             http.StatusBadRequest,
		},
		{
			name:                   "invalid_syntax_code_200",
			err:                    fmt.Errorf("invalid_syntax"),
			code:                   http.StatusOK,
			expectErrorStr:         "invalid_syntax",
			expectHederContentType: HeaderJson,
			expectCode:             http.StatusOK,
		},
		{
			name:                   "price_cannot_be_<=0",
			err:                    fmt.Errorf("price cannot be <= 0"),
			code:                   http.StatusBadRequest,
			expectErrorStr:         "price cannot be <= 0",
			expectHederContentType: HeaderJson,
			expectCode:             http.StatusBadRequest,
		},
		{
			name:                   "price_cannot_be_<=0_code_200",
			err:                    fmt.Errorf("price cannot be <= 0"),
			code:                   http.StatusOK,
			expectErrorStr:         "price cannot be <= 0",
			expectHederContentType: HeaderJson,
			expectCode:             http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			errorResponse(w, tc.err, tc.code)

			if w.Code != tc.expectCode {
				t.Errorf("CODE: got: %d, expect: %d\n", w.Code, tc.expectCode)
			}
			head := w.Header().Get(HeaderContentType)
			if head != tc.expectHederContentType {
				t.Errorf("HEADER: got: %s, expect: %s\n", head, tc.expectHederContentType)
			}

			res := &dto.ErrorDTO{}

			err := json.NewDecoder(w.Body).Decode(res)
			if err != nil {
				t.Fatal(err)
			}

			if res.Error != tc.expectErrorStr {
				t.Errorf("ERROR: got: %s, expect: %s\n", res.Error, tc.expectErrorStr)
			}
			if res.Status != "error" {
				t.Errorf("STATUS: got: %s, expect: error\n", res.Status)
			}
		})
	}
}

func TestCheckHeaderJson(t *testing.T) {
	testCases := []struct {
		name                   string
		valueHeaderContentType string
		expectResult           bool
		expectErrorResponse    bool
	}{
		{
			name:                   "valid_json",
			valueHeaderContentType: HeaderJson,
			expectResult:           true,
			expectErrorResponse:    false,
		},
		{
			name:                   "invalid_wrong_content_type",
			valueHeaderContentType: "application/jsonb",
			expectResult:           false,
			expectErrorResponse:    true,
		},
		{
			name:                   "invalid_empty_content_type",
			valueHeaderContentType: "",
			expectResult:           false,
			expectErrorResponse:    true,
		},
		{
			name:                   "invalid_text_plain",
			valueHeaderContentType: "text/plain",
			expectResult:           false,
			expectErrorResponse:    true,
		},
		{
			name:                   "case_insensitive_test",
			valueHeaderContentType: "APPLICATION/JSON",
			expectResult:           false,
			expectErrorResponse:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/test", nil)
			req.Header.Set(HeaderContentType, tc.valueHeaderContentType)

			result := checkHeaderJson(w, req)

			if result != tc.expectResult {
				t.Errorf("CHECK RESULT: got: %v, expect: %v", result, tc.expectResult)
			}

			if tc.expectErrorResponse {
				if w.Code != http.StatusBadRequest {
					t.Errorf("STATUS CODE: got: %d, expect: %d", w.Code, http.StatusBadRequest)
				}

				contentType := w.Header().Get(HeaderContentType)
				if contentType != HeaderJson {
					t.Errorf("RESPONSE CONTENT-TYPE: got: %s, expect: %s", contentType, HeaderJson)
				}

				var response dto.ErrorDTO
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode error response: %v", err)
				}
				if response.Status != "error" {
					t.Errorf("STATUS: got: %s, expect: error", response.Status)
				}
				if response.Error != "invalid header content-type" {
					t.Errorf("ERROR_STR: got: %s, expect: invalid header content-type", response.Error)
				}
			} else {
				if w.Code != 0 && w.Code != http.StatusOK {
					t.Errorf("Should not write status code, got: %d", w.Code)
				}
				if w.Body.Len() > 0 {
					t.Errorf("Should not write response body, got: %s", w.Body.String())
				}
			}
		})
	}

}
