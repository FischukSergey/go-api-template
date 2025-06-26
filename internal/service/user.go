package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Fisher-Development/woman-app-backend/internal/logger"
	"github.com/Fisher-Development/woman-app-backend/internal/models"
	"github.com/Fisher-Development/woman-app-backend/internal/store"
	"github.com/Fisher-Development/woman-app-backend/internal/types"
	"go.uber.org/zap"
)

// Кастомные ошибки сервиса.
var (
	ErrInvalidUserData = errors.New("invalid user data")
)

// RegistryUser is a service for registering a new user.
type RegistryUser struct {
	storage *store.Storage
}

// IRegistryUser is an interface for registering a new user.
type IRegistryUser interface {
	// CreateUserProfile(ctx context.Context, user *models.User) error
	GetUserByUUID(ctx context.Context, uuid string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
}

// NewRegistryUser creates a new RegistryUser service.
func NewRegistryUser(storage *store.Storage) *RegistryUser {
	return &RegistryUser{storage: storage}
}

// // RegisterUser registers a new user.
// func (r *RegistryUser) RegisterUserProfile(ctx context.Context, user *models.User) error {
// 	// Валидируем данные пользователя
// 	if err := ValidateUserForRegistration(user); err != nil {
// 		// логируем ошибку глобальным логером
// 		logger.GetLogger().Warn("Invalid user data", zap.String("error", err.Error()))
// 		return fmt.Errorf("%w: %v", ErrInvalidUserData, err)
// 	}
// 	if err := r.storage.CreateUserProfile(ctx, user); err != nil {
// 		logger.GetLogger().Warn("Error creating user", zap.String("error", err.Error()))
// 		return err
// 	}
// 	return nil
// }

// UpdateUser is a service for updating a user profile.
func (r *RegistryUser) UpdateUser(ctx context.Context, user *models.User) error {
	// Валидируем данные пользователя
	if err := ValidateUserForUpdate(user); err != nil {
		logger.GetLogger().Warn("Invalid user data", zap.String("error", err.Error()))
		return fmt.Errorf("%w: %v", ErrInvalidUserData, err)
	}
	if err := r.storage.UpdateUser(ctx, user); err != nil {
		logger.GetLogger().Warn("Error updating user", zap.String("error", err.Error()))
		return err
	}
	return nil
}

// UserDashboard is a service for getting a user dashboard.
func (r *RegistryUser) UserDashboard(ctx context.Context, userID types.UserID) (*models.User, error) {
	// получаем пользователя из базы данных
	user, err := r.storage.GetUserByUUID(ctx, userID.String())
	if err != nil {
		logger.GetLogger().Warn("Error getting user dashboard", zap.String("error", err.Error()))
		return nil, err
	}
	return user, nil
}
