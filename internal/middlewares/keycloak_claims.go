package middlewares

import (
	"errors"

	"github.com/Fisher-Development/woman-app-backend/internal/types"
	"github.com/golang-jwt/jwt"
)

var (
	ErrNoAllowedResources = errors.New("no allowed resources")
	ErrSubjectNotDefined  = errors.New(`"sub" is not defined`)
)

type Claims struct {
	jwt.StandardClaims
	// Keycloak использует snake_case для JSON полей
	RealmAccess    map[string][]string `json:"realm_access,omitempty"` //nolint:tagliatelle // Keycloak API format
	ResourceAccess map[string]struct {
		Roles []string `json:"roles,omitempty"`
	} `json:"resource_access,omitempty"` //nolint:tagliatelle // Keycloak API format
}

// Valid returns errors:
// - from StandardClaims validation;
// - ErrNoAllowedResources, if claims doesn't contain `resource_access` map or it's empty;
// - ErrSubjectNotDefined, if claims doesn't contain `sub` field or subject is zero UUID.
func (c Claims) Valid() error {
	// реализуй меня

	if err := c.StandardClaims.Valid(); err != nil {
		return err
	}

	if c.Subject == "" || c.Subject == "00000000-0000-0000-0000-000000000000" {
		return ErrSubjectNotDefined
	}

	if len(c.ResourceAccess) == 0 {
		return ErrNoAllowedResources
	}

	return nil
}

// UserID парсит и возвращает UserID из claims.
func (c Claims) UserID() types.UserID {
	return types.MustParse[types.UserID](c.Subject)
}

// HasResourceRole проверяет наличие указанной роли для указанного ресурса.
func (c Claims) HasResourceRole(resource, role string) bool {
	resourceRoles, exists := c.ResourceAccess[resource]
	if !exists {
		return false
	}

	for _, r := range resourceRoles.Roles {
		if r == role {
			return true
		}
	}
	return false
}
