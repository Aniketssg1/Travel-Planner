package apperrors

import "net/http"

// AppError is a structured application error with an HTTP status code.
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

// --- Factory functions ---

func BadRequest(msg string) *AppError {
	return &AppError{Code: "BAD_REQUEST", Message: msg, HTTPStatus: http.StatusBadRequest}
}

func NotFound(msg string) *AppError {
	return &AppError{Code: "NOT_FOUND", Message: msg, HTTPStatus: http.StatusNotFound}
}

func Conflict(msg string) *AppError {
	return &AppError{Code: "CONFLICT", Message: msg, HTTPStatus: http.StatusConflict}
}

func Unauthorized(msg string) *AppError {
	return &AppError{Code: "UNAUTHORIZED", Message: msg, HTTPStatus: http.StatusUnauthorized}
}

func Forbidden(msg string) *AppError {
	return &AppError{Code: "FORBIDDEN", Message: msg, HTTPStatus: http.StatusForbidden}
}

func Internal(msg string) *AppError {
	return &AppError{Code: "INTERNAL", Message: msg, HTTPStatus: http.StatusInternalServerError}
}

func TooManyRequests(msg string) *AppError {
	return &AppError{Code: "TOO_MANY_REQUESTS", Message: msg, HTTPStatus: http.StatusTooManyRequests}
}
