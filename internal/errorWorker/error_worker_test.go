package error_worker

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"testing"

	logger "github.com/Piccadilly98/subscription_service/internal/loger"
)

func LoadEnv() error {
	err := os.Setenv(EnvNameErrorDb, "connection refused,timeout,network,host not found,authentication failed,password failed,role does not exist,database does not exist,too many connections,SSL required,permission denied,server shutting down")
	if err != nil {
		return err
	}

	err = os.Setenv(EnvNameUserError, "invalid input,validation failed,invalid format,invalid syntax,malformed,input validation failed,invalid data,invalid request,invalid parameters,invalid format for,invalid type for,syntax error,parse error,decoding error,invalid json,malformed json,invalid xml,malformed xml,invalid yaml,invalid encoding,invalid charset,invalid encoding format")
	if err != nil {
		return err
	}
	return nil
}

func TestNewErrorWorker(t *testing.T) {
	err := LoadEnv()
	if err != nil {
		t.Fatalf("fail to load test env: %s\n", err.Error())
	}
	testCases := []struct {
		name             string
		codeDbError      int
		codeUserError    int
		dbLogger         *logger.Logger
		envNameErrorDb   string
		envNameUserError string
		expectedError    error
	}{
		{
			name:             "valid_not_empty_env",
			codeDbError:      http.StatusInternalServerError,
			codeUserError:    http.StatusBadRequest,
			dbLogger:         logger.NewLoggerManager().Get("error db"),
			envNameErrorDb:   EnvNameErrorDb,
			envNameUserError: EnvNameUserError,
		},
		{
			name:             "valid_int_code",
			codeDbError:      500,
			codeUserError:    400,
			dbLogger:         logger.NewLoggerManager().Get("error db"),
			envNameErrorDb:   EnvNameErrorDb,
			envNameUserError: EnvNameUserError,
		},

		{
			name:             "invalid_zero_code_db_error",
			codeDbError:      0,
			codeUserError:    400,
			dbLogger:         logger.NewLoggerManager().Get("error db"),
			envNameErrorDb:   EnvNameErrorDb,
			envNameUserError: EnvNameUserError,
			expectedError:    fmt.Errorf("invalid code DB error"),
		},
		{
			name:             "invalid_zero_code_user_error",
			codeDbError:      500,
			codeUserError:    0,
			dbLogger:         logger.NewLoggerManager().Get("error db"),
			envNameErrorDb:   EnvNameErrorDb,
			envNameUserError: EnvNameUserError,
			expectedError:    fmt.Errorf("invalid code user error"),
		},
		{
			name:             "invalid_negative_code_db_error",
			codeDbError:      -1,
			codeUserError:    400,
			dbLogger:         logger.NewLoggerManager().Get("error db"),
			envNameErrorDb:   EnvNameErrorDb,
			envNameUserError: EnvNameUserError,
			expectedError:    fmt.Errorf("invalid code DB error"),
		},
		{
			name:             "invalid_negative_code_user_error",
			codeDbError:      500,
			codeUserError:    -1,
			dbLogger:         logger.NewLoggerManager().Get("error db"),
			envNameErrorDb:   EnvNameErrorDb,
			envNameUserError: EnvNameUserError,
			expectedError:    fmt.Errorf("invalid code user error"),
		},
		{
			name:             "invalid_loger==nil",
			codeDbError:      500,
			codeUserError:    400,
			dbLogger:         nil,
			envNameErrorDb:   EnvNameErrorDb,
			envNameUserError: EnvNameUserError,
			expectedError:    fmt.Errorf("dbLogger cannot be nil"),
		},

		{
			name:             "invalid_empty_name_db_error",
			codeDbError:      http.StatusInternalServerError,
			codeUserError:    http.StatusBadRequest,
			dbLogger:         logger.NewLoggerManager().Get("error db"),
			envNameErrorDb:   "",
			envNameUserError: EnvNameUserError,
			expectedError:    fmt.Errorf("invalid envNameErrorDb"),
		},
		{
			name:             "invalid_empty_name_user_error",
			codeDbError:      http.StatusInternalServerError,
			codeUserError:    http.StatusBadRequest,
			dbLogger:         logger.NewLoggerManager().Get("error db"),
			envNameErrorDb:   EnvNameErrorDb,
			envNameUserError: "",
			expectedError:    fmt.Errorf("invalid envNameUserError"),
		},

		{
			name:             "invalid_random_name_db_error",
			codeDbError:      http.StatusInternalServerError,
			codeUserError:    http.StatusBadRequest,
			dbLogger:         logger.NewLoggerManager().Get("error db"),
			envNameErrorDb:   "DB_RANDOM_ENV",
			envNameUserError: EnvNameUserError,
			expectedError:    fmt.Errorf("invalid envNameErrorDb"),
		},
		{
			name:             "invalid_random_name_user_error",
			codeDbError:      http.StatusInternalServerError,
			codeUserError:    http.StatusBadRequest,
			dbLogger:         logger.NewLoggerManager().Get("error db"),
			envNameErrorDb:   EnvNameErrorDb,
			envNameUserError: "DB_RANDOM_ENV_USER",
			expectedError:    fmt.Errorf("invalid envNameUserError"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ew, err := NewErrorWorker(tc.codeDbError, tc.codeUserError, tc.dbLogger, tc.envNameErrorDb, tc.envNameUserError)
			if err != nil {
				if tc.expectedError != nil {
					if tc.expectedError.Error() != err.Error() {
						t.Errorf("ERROR: got: %s, expect: %s\n", err.Error(), tc.expectedError.Error())
					}
				} else {
					t.Errorf("unexpected error: %s\n", err.Error())
				}

				if ew != nil {
					t.Errorf("uxepected ew\n")
				}
				return
			} else {
				if tc.expectedError != nil {
					t.Errorf("ERROR: got nil, expect: %s\n", tc.expectedError.Error())
				}
			}

			if ew.codeDbError != tc.codeDbError {
				t.Errorf("CODE_DB_ERROR: got: %d, expect: %d\n", ew.codeDbError, tc.codeDbError)
			}
			if ew.codeUserError != tc.codeUserError {
				t.Errorf("CODE_USER_ERROR: got: %d, expect: %d\n", ew.codeUserError, tc.codeUserError)
			}
		})
	}
}

func TestErrorWorker_CheckErrorGetResult(t *testing.T) {
	err := LoadEnv()
	if err != nil {
		t.Fatalf("fail to load test env: %s\n", err.Error())
	}
	ew, err := NewErrorWorker(http.StatusInternalServerError, http.StatusBadRequest, logger.NewLoggerManager().Get("error"), EnvNameErrorDb, EnvNameUserError)
	if err != nil {
		t.Fatalf("error in new error worker: %s\n", err.Error())
	}
	testCases := []struct {
		name          string
		err           error
		successCode   int
		userErr       string
		expectedCode  int
		expectedError string
	}{
		{
			name:          "custom_error",
			err:           fmt.Errorf("[USER]price cannot be <= 0"),
			successCode:   http.StatusOK,
			userErr:       "invalid_body",
			expectedCode:  http.StatusBadRequest,
			expectedError: "invalid_body: price cannot be <= 0",
		},
		{
			name:          "custom_error_empty",
			err:           fmt.Errorf("[USER]"),
			successCode:   http.StatusOK,
			userErr:       "invalid_body",
			expectedCode:  http.StatusBadRequest,
			expectedError: "invalid_body: ",
		},
		{
			name:          "server_error_connection_refused",
			err:           fmt.Errorf("dial tcp [::1]:5432: connect: connection refused"),
			successCode:   http.StatusOK,
			userErr:       "invalid_body",
			expectedCode:  http.StatusInternalServerError,
			expectedError: "server error",
		},
		{
			name:          "server_error_connection_to_upper",
			err:           fmt.Errorf("DIAL tcp [::1]:5432: connect: connection REFUSED"),
			successCode:   http.StatusOK,
			userErr:       "invalid_body",
			expectedCode:  http.StatusInternalServerError,
			expectedError: "server error",
		},
		{
			name:          "unknown_err",
			err:           fmt.Errorf("random"),
			successCode:   http.StatusOK,
			userErr:       "invalid_body",
			expectedCode:  http.StatusInternalServerError,
			expectedError: "unknown error: random",
		},
		{
			name:          "unknown_err_to_upper",
			err:           fmt.Errorf("RANDOM"),
			successCode:   http.StatusOK,
			userErr:       "invalid_body",
			expectedCode:  http.StatusInternalServerError,
			expectedError: "unknown error: RANDOM",
		},

		{
			name:          "user_error_invalid_input",
			err:           fmt.Errorf("invalid input syntax for type uuid: \"sda\""),
			successCode:   http.StatusOK,
			userErr:       "invalid_body",
			expectedCode:  http.StatusBadRequest,
			expectedError: "invalid_body: invalid input syntax for type uuid: \"sda\"",
		},

		{
			name:          "contex_cancel",
			err:           context.Canceled,
			successCode:   http.StatusOK,
			userErr:       "invalid_body",
			expectedCode:  -1,
			expectedError: "",
		},

		{
			name:          "contex_deadline",
			err:           context.DeadlineExceeded,
			successCode:   http.StatusOK,
			userErr:       "invalid_body",
			expectedCode:  -1,
			expectedError: "",
		},

		{
			name:          "contex_sql_no_rows",
			err:           sql.ErrNoRows,
			successCode:   http.StatusOK,
			userErr:       "invalid_body",
			expectedCode:  http.StatusBadRequest,
			expectedError: "invalid_body: invalid args in request",
		},

		{
			name:          "error_nil",
			err:           nil,
			successCode:   http.StatusOK,
			userErr:       "invalid_body",
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			code, err := ew.CheckErrorGetResult(tc.err, tc.successCode, tc.userErr)
			if code != tc.expectedCode {
				t.Errorf("CODE: got: %d, expect: %d\n", code, tc.expectedCode)
			}
			if err != nil {
				if tc.expectedError != err.Error() {
					t.Errorf("ERROR: got: %s, expect: %s\n", err.Error(), tc.expectedError)
				}
			} else {
				if tc.expectedError != "" {
					t.Errorf("ERROR: got nil, expect: %s\n", tc.expectedError)
				}
			}
		})
	}
}
