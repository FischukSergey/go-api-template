package auth

import (
	"context"
	"net/http"

	"github.com/Fisher-Development/woman-app-backend/api"
	"github.com/go-chi/render"
)

type IAuthService interface {
	RegisterUser(ctx context.Context, req RegisterRequest) error
	LoginUser(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (*LoginResponse, error)
}

// Структура для регистрации.
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"firstName" validate:"omitempty,min=2"`
	LastName  string `json:"lastName" validate:"omitempty,min=2"`
}

// Структура для логина.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Структура для ответа при логине.
type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
}

// Register хендлер для регистрации.
func Register(authService IAuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			api.RespondError(w, r, http.StatusBadRequest, api.ErrorInfo{
				Code:    api.ErrCodeBadRequest,
				Message: err.Error(),
			})
			return
		}

		if err := authService.RegisterUser(r.Context(), req); err != nil {
			// Handle different error types
			api.RespondError(w, r, http.StatusBadRequest, api.ErrorInfo{
				Code:    api.ErrCodeBadRequest,
				Message: err.Error(),
			})
			return
		}

		api.RespondOK(w, r, map[string]string{"status": "registered"})
	}
}

// Login хендлер для логина.
func Login(authService IAuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			api.RespondError(w, r, http.StatusBadRequest, api.ErrorInfo{
				Code:    api.ErrCodeBadRequest,
				Message: err.Error(),
			})
			return
		}

		response, err := authService.LoginUser(r.Context(), req)
		if err != nil {
			api.RespondError(w, r, http.StatusUnauthorized, api.ErrorInfo{
				Code:    api.ErrCodeUnauthorized,
				Message: err.Error(),
			})
			return
		}

		api.RespondOK(w, r, response)
	}
}
