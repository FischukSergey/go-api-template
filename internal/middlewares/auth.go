package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

// Ошибки для извлечения токена.
var (
	ErrMissingAuthHeader = errors.New("missing authorization header")
	ErrInvalidAuthHeader = errors.New("invalid authorization header format")
	ErrEmptyToken        = errors.New("empty token")
)

// AuthMiddleware представляет middleware для авторизации через Keycloak.
type AuthMiddleware struct {
	keycloakClient KeycloakClient
	logger         *zap.Logger
}

// NewAuthMiddleware создает новый экземпляр middleware авторизации.
func NewAuthMiddleware(keycloakClient KeycloakClient) *AuthMiddleware {
	return &AuthMiddleware{
		keycloakClient: keycloakClient,
		logger:         zap.L().Named("auth-middleware"),
	}
}

// RequireAuth возвращает middleware функцию, которая требует валидный JWT токен.
func (a *AuthMiddleware) RequireAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Если Keycloak клиент не настроен, пропускаем проверку
			if a.keycloakClient == nil {
				a.logger.Warn("Keycloak client not configured, skipping auth")
				next.ServeHTTP(w, r)
				return
			}

			// Извлекаем токен из заголовка Authorization
			tokenStr, err := a.extractToken(r)
			if err != nil {
				a.logger.Debug("Failed to extract token", zap.Error(err))
				WriteErrorResponse(w, ErrUnauthorized)
				return
			}

			// Проверяем токен через Keycloak
			tokenResult, err := a.keycloakClient.IntrospectToken(r.Context(), tokenStr)
			if err != nil {
				a.logger.Error("Token introspection failed",
					zap.Error(err),
					zap.String("remote_addr", r.RemoteAddr))
				WriteErrorResponse(w, ErrUnauthorized)
				return
			}

			// Проверяем что токен активен
			if !tokenResult.Active {
				a.logger.Debug("Token is not active",
					zap.String("remote_addr", r.RemoteAddr))
				WriteErrorResponse(w, ErrUnauthorized)
				return
			}

			// Парсим токен, используя наши claims (без проверки подписи, это уже сделал Keycloak)
			tokenClaims := &Claims{}
			token, _ := jwt.ParseWithClaims(tokenStr, tokenClaims, func(_ *jwt.Token) (any, error) {
				return nil, errors.New("не проверяем подпись")
			})

			// Проверяем, что клеймы валидные, включая проверку Subject и ResourceAccess
			if err := tokenClaims.Valid(); err != nil {
				a.logger.Error("Invalid token claims", zap.Error(err))
				WriteErrorResponse(w, ErrUnauthorized)
				return
			}

			// Извлекаем UserID из claims
			userID := tokenClaims.UserID()

			// Добавляем UserID в контекст
			ctx := SetUserID(r.Context(), userID)
			ctx = SetToken(ctx, tokenStr)

			// Добавляем роли если есть
			if realmRoles, exists := tokenClaims.RealmAccess["roles"]; exists && len(realmRoles) > 0 {
				ctx = SetUserRoles(ctx, realmRoles)
			}

			// Сохраняем токен в контекст
			if token == nil {
				token = &jwt.Token{Claims: tokenClaims}
			}
			ctx = SetJWTToken(ctx, token)

			a.logger.Debug("User authenticated successfully",
				zap.String("user_id", userID.String()),
				zap.String("remote_addr", r.RemoteAddr))

			// Передаем управление следующему обработчику
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extractToken извлекает JWT токен из заголовка Authorization.
func (a *AuthMiddleware) extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrMissingAuthHeader
	}

	// Проверяем формат "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", ErrInvalidAuthHeader
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", ErrEmptyToken
	}

	return token, nil
}
