package error_worker

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	logger "github.com/Piccadilly98/subscription_service/internal/loger"
)

const (
	LevelForDBError   = "CRITICAL"
	LevelForUserError = "WARNING"

	EnvNameErrorDb   = "DB_CONNECTION_ERROR"
	EnvNameUserError = "DB_USER_ERROR"
)

type ErrorWorker struct {
	dbError       []string
	codeDbError   int
	userError     []string
	codeUserError int
	dbLogger      *logger.Logger
}

func NewErrorWorker(
	codeDbError int,
	codeUserError int,
	dbLogger *logger.Logger,
) *ErrorWorker {
	ew := &ErrorWorker{
		codeDbError:   codeDbError,
		codeUserError: codeUserError,
		dbLogger:      dbLogger,
	}

	errDb := os.Getenv(EnvNameErrorDb)
	errUser := os.Getenv(EnvNameUserError)
	ew.dbError = strings.Split(errDb, ",")
	ew.userError = strings.Split(errUser, ",")
	return ew
}

func (e *ErrorWorker) CheckErrorGetResult(err error, successCode int, userErr string) (int, error) {
	if err == nil {
		return successCode, nil
	}
	errStr := err.Error()
	if isContextError(err) {
		return -1, nil
	}
	if err == sql.ErrNoRows {
		e.dbLogger.Printf(LevelForUserError, "User error: %s", errStr)
		return e.codeUserError, fmt.Errorf("%s: invalid args in request", userErr)
	}

	if strings.HasPrefix(errStr, "[USER]") {
		cleanErr := strings.TrimPrefix(errStr, "[USER]")
		return e.codeUserError, fmt.Errorf("%s: %s", userErr, cleanErr)
	}

	for _, pattern := range e.dbError {
		if strings.Contains(strings.ToLower(errStr), strings.ToLower(pattern)) {
			e.dbLogger.Printf(LevelForDBError, "DB error: %s", errStr)
			return e.codeDbError, fmt.Errorf("server error")
		}
	}

	for _, pattern := range e.userError {
		if strings.Contains(strings.ToLower(errStr), strings.ToLower(pattern)) {
			e.dbLogger.Printf(LevelForUserError, "User error: %s", errStr)
			return e.codeUserError, fmt.Errorf("%s: %s", userErr, errStr)
		}
	}
	e.dbLogger.Printf(LevelForDBError, "Unknown error: %s", errStr)
	return http.StatusInternalServerError, fmt.Errorf("unknown error: %s", errStr)
}

func isContextError(err error) bool {
	return errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, context.Canceled)
}
