package api

import (
	"net/http"

	"github.com/go-chi/render"
)

const (
	ErrCodeValidationFailed = "VALIDATION_FAILED"
	ErrCodeInternalServer   = "INTERNAL_SERVER_ERROR"
	ErrCodeNotFound         = "NOT_FOUND"
	ErrCodeConflict         = "CONFLICT"
	ErrCodeUnauthorized     = "UNAUTHORIZED"
	ErrCodeBadRequest       = "BAD_REQUEST"
)

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// RespondOK отправляет успешный ответ клиенту.
func RespondOK(w http.ResponseWriter, r *http.Request, data any) {
	w.Header().Set("Content-Type", "application/json")
	render.Status(r, http.StatusOK)
	render.JSON(w, r, data)
}

// RespondError отправляет ошибку клиенту.
func RespondError(w http.ResponseWriter, r *http.Request, status int, err ErrorInfo) {
	w.Header().Set("Content-Type", "application/json")
	render.Status(r, status)
	render.JSON(w, r, ErrorInfo{
		Code:    err.Code,
		Message: err.Message,
	})
}
