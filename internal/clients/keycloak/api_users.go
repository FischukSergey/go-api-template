package keycloakclient

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Fisher-Development/woman-app-backend/internal/logger"
	"github.com/Fisher-Development/woman-app-backend/internal/types"
	"go.uber.org/zap"
)

// CreateUserRequest структура для создания пользователя.
type CreateUserRequest struct {
	Username      string                         `json:"username"`
	Email         string                         `json:"email"`
	FirstName     string                         `json:"firstName"`
	LastName      string                         `json:"lastName"`
	Enabled       bool                           `json:"enabled"`
	EmailVerified bool                           `json:"emailVerified"`
	Credentials   []UserCredentialRepresentation `json:"credentials,omitempty"`
}

// UserCredentialRepresentation для установки пароля.
type UserCredentialRepresentation struct {
	Type      string `json:"type"`
	Value     string `json:"value"`
	Temporary bool   `json:"temporary"`
}

// CreateUserResponse ответ после создания пользователя.
type CreateUserResponse struct {
	UserID    types.UserID `json:"user_id"`
	Email     string       `json:"email"`
	FirstName string       `json:"firstName"`
	LastName  string       `json:"lastName"`
}

// TokenResponse структура ответа при получении токена.
type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

// CreateUser создает нового пользователя в Keycloak.
func (c *Client) CreateUser(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error) {
	// Сначала получаем admin токен
	adminToken, err := c.getAdminToken(ctx)
	if err != nil {
		logger.GetLogger().Error("Failed to get admin token", zap.Error(err))
		return nil, fmt.Errorf("get admin token: %w", err)
	}

	url := fmt.Sprintf("%s/admin/realms/%s/users", c.basePath, c.realm)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetAuthToken(adminToken). // Используем Bearer токен
		SetBody(req).
		Post(url)
	if err != nil {
		logger.GetLogger().Error("Failed to create user",
			zap.String("url", url),
			zap.Error(err))
		return nil, fmt.Errorf("create user request failed: %w", err)
	}

	// Проверяем статус код
	if resp.StatusCode() != http.StatusCreated {
		logger.GetLogger().Error("Create user failed with status",
			zap.Int("status_code", resp.StatusCode()),
			zap.String("response", string(resp.Body())))
		return nil, fmt.Errorf("create user failed with status %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	// Keycloak возвращает Location header с ID пользователя
	location := resp.Header().Get("Location")
	if location == "" {
		logger.GetLogger().Error("No Location header in create user response")
		return nil, errors.New("no Location header in response")
	}

	// Извлекаем ID из Location header
	// Location: http://localhost:8080/admin/realms/myrealm/users/12345678-1234-1234-1234-123456789012
	parts := strings.Split(location, "/")
	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid Location header format: %s", location)
	}

	userIDStr := parts[len(parts)-1]

	// Используем ваш кастомный Parse для types.UserID
	userID, err := types.Parse[types.UserID](userIDStr)
	if err != nil {
		logger.GetLogger().Error("Failed to parse user ID from Location header",
			zap.String("location", location),
			zap.String("user_id_str", userIDStr),
			zap.Error(err))
		return nil, fmt.Errorf("failed to parse user ID: %w", err)
	}

	logger.GetLogger().Info("User created successfully",
		zap.String("user_id", userID.String()),
		zap.String("username", req.Username))

	return &CreateUserResponse{
		UserID:    userID,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}, nil
}

// LoginUser аутентифицирует пользователя.
func (c *Client) LoginUser(ctx context.Context, email, password string) (*TokenResponse, error) {
	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", c.basePath, c.realm)

	var tokenResp TokenResponse

	resp, err := c.cli.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"grant_type":    "password",
			"username":      email,
			"password":      password,
			"client_id":     c.clientID,
			"client_secret": c.clientSecret,
		}).
		SetResult(&tokenResp).
		Post(url)
	if err != nil {
		return nil, fmt.Errorf("send login request: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("login failed: %v - %s", resp.Status(), resp.String())
	}

	return &tokenResp, nil
}

// RefreshAccessToken обновляет access token используя refresh token.
func (c *Client) RefreshAccessToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", c.basePath, c.realm)

	var tokenResp TokenResponse

	resp, err := c.cli.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"grant_type":    "refresh_token",
			"refresh_token": refreshToken,
			"client_id":     c.clientID,
			"client_secret": c.clientSecret,
		}).
		SetResult(&tokenResp).
		Post(url)
	if err != nil {
		return nil, fmt.Errorf("refresh token request: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("refresh failed: %v - %s", resp.Status(), resp.String())
	}

	return &tokenResp, nil
}

// getAdminToken получает токен для Admin API.
func (c *Client) getAdminToken(ctx context.Context) (string, error) {
	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", c.basePath, c.realm)

	type AdminTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	var tokenResp AdminTokenResponse

	resp, err := c.cli.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     c.clientID,
			"client_secret": c.clientSecret,
		}).
		SetResult(&tokenResp).
		Post(url)
	if err != nil {
		return "", fmt.Errorf("send token request: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		logger.GetLogger().Error("Failed to get admin token",
			zap.Int("status_code", resp.StatusCode()),
			zap.String("response", string(resp.Body())))
		return "", fmt.Errorf("token request failed: %v - %s", resp.Status(), resp.String())
	}

	return tokenResp.AccessToken, nil
}
