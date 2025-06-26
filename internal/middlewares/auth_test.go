package middlewares_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	keycloakclient "github.com/Fisher-Development/woman-app-backend/internal/clients/keycloak"
	"github.com/Fisher-Development/woman-app-backend/internal/middlewares"
	middlewaresmocks "github.com/Fisher-Development/woman-app-backend/internal/middlewares/mocks"
	"github.com/Fisher-Development/woman-app-backend/internal/types"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const bearerPrefix = "Bearer "

func TestAuthMiddleware(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareSuite))
}

type AuthMiddlewareSuite struct {
	suite.Suite
	ctrl           *gomock.Controller
	keycloakClient *middlewaresmocks.MockKeycloakClient
	authMdlwr      *middlewares.AuthMiddleware
	router         *chi.Mux
}

func (s *AuthMiddlewareSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.keycloakClient = middlewaresmocks.NewMockKeycloakClient(s.ctrl)
	s.authMdlwr = middlewares.NewAuthMiddleware(s.keycloakClient)

	s.router = chi.NewRouter()
	s.router.Use(s.authMdlwr.RequireAuth())
	s.router.Get("/test", s.testHandler)
}

func (s *AuthMiddlewareSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *AuthMiddlewareSuite) testHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := middlewares.GetUserID(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"user_id":"` + userID.String() + `"}`))
}

// Положительные тесты

func (s *AuthMiddlewareSuite) TestValidToken() {
	//nolint:lll,gosec // Test JWT token for middleware testing
	const token = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJIR1lJcHN1UXlsZFNJZTB1T0JaeEpuQjBkZlFuTWI5LUlFcmx6NHk5ek9BIn0.eyJleHAiOjI2NjcxOTk1ODAsImlhdCI6MTY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiNWNiNDBkYzAtYTI0OS00NzgzLWEzMDEtOWUxZjNjZjNlYTQxIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiY2hhdC11aS1jbGllbnQiLCJub25jZSI6ImJhMzdmZDVhLThjMzktNDgxNC1hZmNiLTk1MmExOGI3MjY3ZCIsInNlc3Npb25fc3RhdGUiOiJkODZkMTk4ZS1jMWM1LTRlZGQtODM1MC0zNjFlZTU4MTcxZjIiLCJhY3IiOiIwIiwiYWxsb3dlZC1vcmlnaW5zIjpbIiIsIioiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbIm9mZmxpbmVfYWNjZXNzIiwiZGVmYXVsdC1yb2xlcy1iYW5rIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInJlc291cmNlX2FjY2VzcyI6eyJjaGF0LXVpLWNsaWVudCI6eyJyb2xlcyI6WyJzdXBwb3J0LWNoYXQtY2xpZW50Il19LCJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50IiwibWFuYWdlLWFjY291bnQtbGlua3MiLCJ2aWV3LXByb2ZpbGUiXX19LCJzY29wZSI6Im9wZW5pZCBwcm9maWxlIGVtYWlsIiwic2lkIjoiZDg2ZDE5OGUtYzFjNS00ZWRkLTgzNTAtMzYxZWU1ODE3MWYyIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsInByZWZlcnJlZF91c2VybmFtZSI6ImJvbmQwMDciLCJnaXZlbl9uYW1lIjoiIiwiZmFtaWx5X25hbWUiOiIiLCJlbWFpbCI6ImJvbmQwMDdAdWsuY29tIn0.we-dont-check-signature"

	s.keycloakClient.EXPECT().
		IntrospectToken(gomock.Any(), token).
		Return(&keycloakclient.IntrospectTokenResult{Active: true}, nil)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", bearerPrefix+token)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "5cb40dc0-a249-4783-a301-9e1f3cf3ea41")
}

// Отрицательные тесты

func (s *AuthMiddlewareSuite) TestNoAuthorizationHeader() {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusUnauthorized, w.Code)
	s.Contains(w.Body.String(), "unauthorized")
}

func (s *AuthMiddlewareSuite) TestInvalidAuthHeader() {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Basic token123")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusUnauthorized, w.Code)
	s.Contains(w.Body.String(), "unauthorized")
}

func (s *AuthMiddlewareSuite) TestEmptyToken() {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer ")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *AuthMiddlewareSuite) TestKeycloakIntrospectError() {
	const token = "valid.jwt.token" //nolint:gosec // Test token string

	s.keycloakClient.EXPECT().
		IntrospectToken(gomock.Any(), token).
		Return(nil, context.Canceled)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", bearerPrefix+token)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *AuthMiddlewareSuite) TestInactiveToken() {
	const token = "inactive.jwt.token" //nolint:gosec // Test token string

	s.keycloakClient.EXPECT().
		IntrospectToken(gomock.Any(), token).
		Return(&keycloakclient.IntrospectTokenResult{Active: false}, nil)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", bearerPrefix+token)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *AuthMiddlewareSuite) TestInvalidTokenClaims() {
	//nolint:lll,gosec // Test JWT token for middleware testing
	const expiredToken = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJIR1lJcHN1UXlsZFNJZTB1T0JaeEpuQjBkZlFuTWI5LUlFcmx6NHk5ek9BIn0.eyJleHAiOjE2NjcxOTk1ODAsImlhdCI6MTY2NzE5OTI4MH0.we-dont-check-signature"

	s.keycloakClient.EXPECT().
		IntrospectToken(gomock.Any(), expiredToken).
		Return(&keycloakclient.IntrospectTokenResult{Active: true}, nil)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", bearerPrefix+expiredToken)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *AuthMiddlewareSuite) TestNoKeycloakClient() {
	// Создаем middleware без Keycloak клиента
	authMdlwr := middlewares.NewAuthMiddleware(nil)
	router := chi.NewRouter()
	router.Use(authMdlwr.RequireAuth())
	router.Get("/test", s.testHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Должен пропустить без проверки
	s.Equal(http.StatusInternalServerError, w.Code) // Поскольку userID не найден
}

// Тест функций работы с контекстом

func TestContextFunctions(t *testing.T) {
	ctx := context.Background()
	userID := types.MustParse[types.UserID]("5cb40dc0-a249-4783-a301-9e1f3cf3ea41")

	// Тестируем SetUserID/GetUserID
	ctx = middlewares.SetUserID(ctx, userID)
	retrievedUserID, err := middlewares.GetUserID(ctx)
	require.NoError(t, err)
	assert.Equal(t, userID, retrievedUserID)

	// Тестируем MustGetUserID
	mustUserID := middlewares.MustGetUserID(ctx)
	assert.Equal(t, userID, mustUserID)

	// Тестируем GetUserID с пустым контекстом
	emptyCtx := context.Background()
	_, err = middlewares.GetUserID(emptyCtx)
	require.Error(t, err)
	assert.Equal(t, middlewares.ErrUserIDNotFound, err)

	// Тестируем роли
	roles := []string{"admin", "user"}
	ctx = middlewares.SetUserRoles(ctx, roles)
	retrievedRoles, ok := middlewares.GetUserRoles(ctx)
	assert.True(t, ok)
	assert.Equal(t, roles, retrievedRoles)

	// Тестируем токен
	token := "test.jwt.token" //nolint:gosec // Test token string
	ctx = middlewares.SetToken(ctx, token)
	retrievedToken, ok := middlewares.GetToken(ctx)
	assert.True(t, ok)
	assert.Equal(t, token, retrievedToken)
}

func TestMustGetUserID_Panic(t *testing.T) {
	ctx := context.Background()

	assert.Panics(t, func() {
		middlewares.MustGetUserID(ctx)
	})
}
