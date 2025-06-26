package user

import (
	"context"
	"errors"
	"net/http"

	"github.com/Fisher-Development/woman-app-backend/api"
	"github.com/Fisher-Development/woman-app-backend/internal/middlewares"
	"github.com/Fisher-Development/woman-app-backend/internal/models"
	"github.com/Fisher-Development/woman-app-backend/internal/service"
	"github.com/Fisher-Development/woman-app-backend/internal/store"
	"github.com/Fisher-Development/woman-app-backend/internal/types"
	"github.com/go-chi/render"
)

// IRegistryUser is an interface for registering a new user.
type IRegistryUser interface {
	// RegisterUserProfile(ctx context.Context, user *models.User) error
	UserDashboard(ctx context.Context, userID types.UserID) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
}

// // UserRegistry is a handler for registering a new user.
// func RegistryProfile(registry IRegistryUser) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var user models.User
// 		// получаем UUID из контекста
// 		userID, ok := middlewares.GetUserFromContext(r.Context())
// 		if !ok {
// 			api.RespondError(w, r, http.StatusUnauthorized, api.ErrorInfo{
// 				Code:    api.ErrCodeUnauthorized,
// 				Message: "User not found keycloak",
// 			})
// 			return
// 		}
// 		// декодируем тело запроса в структуру user
// 		if err := render.DecodeJSON(r.Body, &user); err != nil {
// 			api.RespondError(w, r, http.StatusBadRequest, api.ErrorInfo{
// 				Code:    api.ErrCodeBadRequest,
// 				Message: err.Error(),
// 			})
// 			return
// 		}
// 		user.UUID = userID.String()

// 		// регистрируем профайл пользователя (валидация в сервисном слое)
// 		err := registry.RegisterUserProfile(r.Context(), &user)
// 		if err != nil {
// 			switch err {
// 			// Если ошибка валидации - возвращаем 400.
// 			case service.ErrInvalidUserData:
// 				api.RespondError(w, r, http.StatusBadRequest, api.ErrorInfo{
// 					Code:    api.ErrCodeBadRequest,
// 					Message: err.Error(),
// 				})
// 			// Если пользователь уже существует - возвращаем 409
// 			case store.ErrUserAlreadyExists:
// 				api.RespondError(w, r, http.StatusConflict, api.ErrorInfo{
// 					Code:    api.ErrCodeConflict,
// 					Message: err.Error(),
// 				})
// 			// Если ошибка неизвестна - возвращаем 500
// 			default:
// 				api.RespondError(w, r, http.StatusInternalServerError, api.ErrorInfo{
// 					Code:    api.ErrCodeInternalServer,
// 					Message: err.Error(),
// 				})
// 			}
// 			return
// 		}
// 		// отправляем ответ клиенту
// 		api.RespondOK(w, r, map[string]string{"status": "ok"})
// 	}
// }

// Update is a handler for updating a user profile.
func Update(registry IRegistryUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		// получаем UUID из контекста
		userID, ok := middlewares.GetUserFromContext(r.Context())
		if !ok {
			api.RespondError(w, r, http.StatusUnauthorized, api.ErrorInfo{
				Code:    api.ErrCodeUnauthorized,
				Message: "User not found keycloak",
			})
			return
		}
		// декодируем тело запроса в структуру user
		if err := render.DecodeJSON(r.Body, &user); err != nil {
			api.RespondError(w, r, http.StatusBadRequest, api.ErrorInfo{
				Code:    api.ErrCodeBadRequest,
				Message: err.Error(),
			})
			return
		}
		user.UUID = userID.String()
		err := registry.UpdateUser(r.Context(), &user)
		if err != nil {
			if errors.Is(err, service.ErrInvalidUserData) {
				api.RespondError(w, r, http.StatusBadRequest, api.ErrorInfo{
					Code:    api.ErrCodeBadRequest,
					Message: err.Error(),
				})
				return
			}
			if errors.Is(err, store.ErrUserNotFound) {
				api.RespondError(w, r, http.StatusNotFound, api.ErrorInfo{
					Code:    api.ErrCodeNotFound,
					Message: err.Error(),
				})
				return
			}
			if errors.Is(err, store.ErrEmailAlreadyExists) {
				api.RespondError(w, r, http.StatusConflict, api.ErrorInfo{
					Code:    api.ErrCodeConflict,
					Message: err.Error(),
				})
				return
			}
			// default case
			api.RespondError(w, r, http.StatusInternalServerError, api.ErrorInfo{
				Code:    api.ErrCodeInternalServer,
				Message: err.Error(),
			})
			return
		}
		api.RespondOK(w, r, map[string]string{"status": "ok"})
	}
}

// Dashboard is a handler for getting a user dashboard.
func Dashboard(registry IRegistryUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// получаем UUID из контекста
		userID, ok := middlewares.GetUserFromContext(r.Context())
		if !ok {
			api.RespondError(w, r, http.StatusUnauthorized, api.ErrorInfo{
				Code:    api.ErrCodeUnauthorized,
				Message: "User not found keycloak",
			})
			return
		}
		// получаем пользователя из базы данных
		user, err := registry.UserDashboard(r.Context(), userID)
		if err != nil {
			if errors.Is(err, store.ErrUserNotFound) {
				api.RespondError(w, r, http.StatusNotFound, api.ErrorInfo{
					Code:    api.ErrCodeNotFound,
					Message: err.Error(),
				})
				return
			}
			// default case
			api.RespondError(w, r, http.StatusInternalServerError, api.ErrorInfo{
				Code:    api.ErrCodeInternalServer,
				Message: "Internal server error",
			})
			return
		}
		// отправляем ответ клиенту
		api.RespondOK(w, r, user)
	}
}
