package errors

import (
	"encoding/json"
	"net/http"
)

type ErrorCode string

const (
	ErrCodeInvalidRequest    ErrorCode = "INVALID_REQUEST"
	ErrCodeInvalidEmail      ErrorCode = "INVALID_EMAIL"
	ErrCodeInvalidPassword   ErrorCode = "INVALID_PASSWORD"
	ErrCodeUnauthorized      ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden         ErrorCode = "FORBIDDEN"
	ErrCodeNotFound          ErrorCode = "NOT_FOUND"
	ErrCodeConflict          ErrorCode = "CONFLICT"
	ErrCodeValidationFailed  ErrorCode = "VALIDATION_FAILED"
	ErrCodeInternalError     ErrorCode = "INTERNAL_ERROR"
	ErrCodeDatabaseError     ErrorCode = "DATABASE_ERROR"
	ErrCodeInsufficientStock ErrorCode = "INSUFFICIENT_STOCK"
)

type APIError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details any       `json:"details,omitempty"`
}

func (e APIError) Error() string {
	return e.Message
}

func New(code ErrorCode, message string) APIError {
	return APIError{
		Code:    code,
		Message: message,
	}
}

func NewWithDetails(code ErrorCode, message string, details any) APIError {
	return APIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func RespondWithError(w http.ResponseWriter, statusCode int, err APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(err)
}

func RespondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}
