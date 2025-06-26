package auth

import (
	"net/http"

	"github.com/Fisher-Development/woman-app-backend/api"
	"github.com/go-chi/render"
)

// RefreshToken.
func RefreshToken(authService IAuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type RefreshRequest struct {
			RefreshToken string `json:"refreshToken"`
		}

		var req RefreshRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			api.RespondError(w, r, http.StatusBadRequest, api.ErrorInfo{
				Code:    api.ErrCodeBadRequest,
				Message: "Invalid request",
			})
			return
		}

		// Обновляем токен через Keycloak
		tokens, err := authService.RefreshAccessToken(r.Context(), req.RefreshToken)
		if err != nil {
			api.RespondError(w, r, http.StatusUnauthorized, api.ErrorInfo{
				Code:    api.ErrCodeUnauthorized,
				Message: "Invalid refresh token",
			})
			return
		}

		api.RespondOK(w, r, tokens)
	}
}
