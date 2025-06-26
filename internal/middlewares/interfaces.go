package middlewares

import (
	"context"

	keycloakclient "github.com/Fisher-Development/woman-app-backend/internal/clients/keycloak"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/keycloak_client_mock.go

type KeycloakClient interface {
	IntrospectToken(ctx context.Context, token string) (*keycloakclient.IntrospectTokenResult, error)
}
