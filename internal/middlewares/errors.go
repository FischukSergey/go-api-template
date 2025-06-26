package middlewares

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse структура ответа об ошибке.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Предопределенные ошибки.
var (
	ErrUnauthorized = &ErrorResponse{
		Error:   "unauthorized",
		Message: "Authentication required",
		Code:    http.StatusUnauthorized,
	}

	ErrForbidden = &ErrorResponse{
		Error:   "forbidden",
		Message: "Access denied",
		Code:    http.StatusForbidden,
	}

	ErrInternalError = &ErrorResponse{
		Error:   "internal_error",
		Message: "Internal server error",
		Code:    http.StatusInternalServerError,
	}
)

// WriteErrorResponse отправляет JSON ответ с ошибкой.
func WriteErrorResponse(w http.ResponseWriter, err *ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)

	if encodeErr := json.NewEncoder(w).Encode(err); encodeErr != nil {
		// Если не удалось закодировать JSON, отправляем простой текст
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Server Error"))
	}
}
