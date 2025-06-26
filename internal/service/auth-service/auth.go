package authservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/Fisher-Development/woman-app-backend/api/auth"
	keycloakclient "github.com/Fisher-Development/woman-app-backend/internal/clients/keycloak"
	"github.com/Fisher-Development/woman-app-backend/internal/models"
	"github.com/Fisher-Development/woman-app-backend/internal/store"
	"go.uber.org/zap"
)

// AuthService реализует IAuthService интерфейс.
type AuthService struct {
	storage             *store.Storage
	keycloakClient      *keycloakclient.Client // для аутентификации
	keycloakAdminClient *keycloakclient.Client // для создания пользователей
}

// NewAuthService создает новый AuthService.
func NewAuthService(
	storage *store.Storage,
	keycloakClient *keycloakclient.Client,
	keycloakAdminClient *keycloakclient.Client,
) *AuthService {
	return &AuthService{
		storage:             storage,
		keycloakClient:      keycloakClient,
		keycloakAdminClient: keycloakAdminClient,
	}
}

// RegisterUser регистрирует нового пользователя.
func (s *AuthService) RegisterUser(ctx context.Context, req auth.RegisterRequest) error {
	logger := zap.L().Named("auth-service")

	logger.Info("Starting user registration",
		zap.String("email", req.Email),
		zap.String("firstName", req.FirstName))

	// Используем email как username если не указан отдельно
	username := req.Email
	if username == "" {
		return errors.New("email is required for registration")
	}

	// Создаем запрос для Keycloak
	createReq := keycloakclient.CreateUserRequest{
		Username:      username,
		Email:         req.Email,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Enabled:       true,
		EmailVerified: false, // потребует верификации email
		Credentials: []keycloakclient.UserCredentialRepresentation{
			{
				Type:      "password",
				Value:     req.Password,
				Temporary: false,
			},
		},
	}

	// Создаем пользователя в Keycloak
	keycloakUser, err := s.keycloakAdminClient.CreateUser(ctx, createReq)
	if err != nil {
		logger.Error("Failed to create user in Keycloak", zap.Error(err))
		return fmt.Errorf("failed to create user in Keycloak: %v", err)
	}

	logger.Info("User created in Keycloak",
		zap.String("user_id", keycloakUser.UserID.String()))

	// Создаем пользователя в нашей БД
	user := &models.User{
		UUID:      keycloakUser.UserID.String(), // используем ID из Keycloak
		Email:     keycloakUser.Email,
		FirstName: keycloakUser.FirstName,
		LastName:  keycloakUser.LastName,
	}

	if err := s.storage.CreateUser(ctx, user); err != nil {
		logger.Error("Failed to create user in database", zap.Error(err))
		// СДЕЛАТЬ В идеале здесь нужно откатить создание в Keycloak
		return fmt.Errorf("failed to create user in database: %v", err)
	}

	logger.Info("User registration completed successfully",
		zap.String("user_id", keycloakUser.UserID.String()))

	return nil
}

// LoginUser аутентифицирует пользователя.
func (s *AuthService) LoginUser(ctx context.Context, req auth.LoginRequest) (*auth.LoginResponse, error) {
	logger := zap.L().Named("auth-service")

	logger.Info("User login attempt", zap.String("email", req.Email))

	// Аутентифицируем через Keycloak
	tokenResp, err := s.keycloakClient.LoginUser(ctx, req.Email, req.Password)
	if err != nil {
		logger.Warn("Login failed",
			zap.String("email", req.Email),
			zap.Error(err))
		return nil, fmt.Errorf("authentication failed: %v", err)
	}

	logger.Info("User logged in successfully", zap.String("email", req.Email))

	return &auth.LoginResponse{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresIn:    tokenResp.ExpiresIn,
	}, nil
}

// RefreshAccessToken обновляет access token.
func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (*auth.LoginResponse, error) {
	logger := zap.L().Named("auth-service")

	logger.Info("Refreshing access token")

	tokenResp, err := s.keycloakClient.RefreshAccessToken(ctx, refreshToken)
	if err != nil {
		logger.Warn("Token refresh failed", zap.Error(err))
		return nil, fmt.Errorf("token refresh failed: %v", err)
	}

	logger.Info("Access token refreshed successfully")

	return &auth.LoginResponse{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresIn:    tokenResp.ExpiresIn,
	}, nil
}
