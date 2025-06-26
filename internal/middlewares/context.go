package middlewares

import (
	"context"
	"errors"

	keycloakclient "github.com/Fisher-Development/woman-app-backend/internal/clients/keycloak"
	"github.com/Fisher-Development/woman-app-backend/internal/types"
	"github.com/golang-jwt/jwt"
)

// ContextKey тип для ключей контекста.
type ContextKey string

const (
	// UserIDKey ключ для хранения UserID в контексте.
	UserIDKey ContextKey = "user_id"
	// UserRolesKey ключ для хранения ролей пользователя.
	UserRolesKey ContextKey = "user_roles"
	// TokenKey ключ для хранения токена.
	TokenKey ContextKey = "token"
	// JWTTokenKey ключ для хранения JWT токена.
	JWTTokenKey ContextKey = "jwt_token"
)

// Ошибки извлечения из контекста.
var (
	ErrUserIDNotFound = errors.New("user ID not found in context")
	ErrUserIDInvalid  = errors.New("user ID has invalid type")
)

// GetUserFromContext извлекает информацию о пользователе из контекста.
func GetUserFromContext(ctx context.Context) (types.UserID, bool) {
	user, ok := ctx.Value(UserIDKey).(types.UserID)
	return user, ok
}

// GetTokenFromContext извлекает информацию о токене из контекста.
func GetTokenFromContext(ctx context.Context) (*keycloakclient.IntrospectTokenResult, bool) {
	token, ok := ctx.Value(TokenKey).(*keycloakclient.IntrospectTokenResult)
	return token, ok
}

// IsAuthenticated проверяет аутентифицирован ли пользователь.
func IsAuthenticated(ctx context.Context) bool {
	user, ok := GetUserFromContext(ctx)
	return ok && !user.IsZero()
}

// MustGetUser извлекает пользователя из контекста, паникует если пользователь не найден.
// Используйте только в защищенных обработчиках после RequireAuth middleware.
func MustGetUser(ctx context.Context) types.UserID {
	user, ok := GetUserFromContext(ctx)
	if !ok {
		panic("user not found in context - ensure RequireAuth middleware is used")
	}
	return user
}

// SetUserID добавляет UserID в контекст.
func SetUserID(ctx context.Context, userID types.UserID) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// GetUserID извлекает UserID из контекста.
func GetUserID(ctx context.Context) (types.UserID, error) {
	value := ctx.Value(UserIDKey)
	if value == nil {
		return types.UserIDNil, ErrUserIDNotFound
	}

	userID, ok := value.(types.UserID)
	if !ok {
		return types.UserIDNil, ErrUserIDInvalid
	}

	return userID, nil
}

// MustGetUserID извлекает UserID из контекста или паникует.
func MustGetUserID(ctx context.Context) types.UserID {
	userID, err := GetUserID(ctx)
	if err != nil {
		panic(err)
	}
	return userID
}

// SetUserRoles добавляет роли пользователя в контекст.
func SetUserRoles(ctx context.Context, roles []string) context.Context {
	return context.WithValue(ctx, UserRolesKey, roles)
}

// GetUserRoles извлекает роли пользователя из контекста.
func GetUserRoles(ctx context.Context) ([]string, bool) {
	value := ctx.Value(UserRolesKey)
	if value == nil {
		return nil, false
	}

	roles, ok := value.([]string)
	return roles, ok
}

// SetToken добавляет токен в контекст.
func SetToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, TokenKey, token)
}

// GetToken извлекает токен из контекста.
func GetToken(ctx context.Context) (string, bool) {
	value := ctx.Value(TokenKey)
	if value == nil {
		return "", false
	}

	token, ok := value.(string)
	return token, ok
}

// SetJWTToken добавляет JWT токен в контекст.
func SetJWTToken(ctx context.Context, token *jwt.Token) context.Context {
	return context.WithValue(ctx, JWTTokenKey, token)
}

// GetJWTToken извлекает JWT токен из контекста.
func GetJWTToken(ctx context.Context) (*jwt.Token, bool) {
	value := ctx.Value(JWTTokenKey)
	if value == nil {
		return nil, false
	}

	token, ok := value.(*jwt.Token)
	return token, ok
}

// GetClaims извлекает claims из контекста.
func GetClaims(ctx context.Context) (*Claims, bool) {
	token, ok := GetJWTToken(ctx)
	if !ok {
		return nil, false
	}

	tokenClaims, ok := token.Claims.(*Claims)
	return tokenClaims, ok
}
