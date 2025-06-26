package service_test

import (
	"strings"
	"testing"

	"github.com/Fisher-Development/woman-app-backend/internal/models"
	"github.com/Fisher-Development/woman-app-backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test_validateUserForRegistration тестирует основную функцию валидации.
func Test_validateUserForRegistration(t *testing.T) {
	tests := []struct {
		name    string
		user    *models.User
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid user - all required fields",
			user: &models.User{
				UUID:  "123e4567-e89b-12d3-a456-426614174000",
				Email: "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "valid user - with optional fields",
			user: &models.User{
				UUID:      "123e4567-e89b-12d3-a456-426614174000",
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  "Doe",
				Sex:       "male",
				City:      "Moscow",
				Country:   "Russia",
			},
			wantErr: false,
		},
		{
			name: "invalid - missing UUID",
			user: &models.User{
				Email:     "test@example.com",
				FirstName: "John",
			},
			wantErr: true,
			errMsg:  "user UUID is required",
		},
		{
			name: "invalid - missing email",
			user: &models.User{
				UUID:      "123e4567-e89b-12d3-a456-426614174000",
				FirstName: "John",
			},
			wantErr: true,
			errMsg:  "email is required",
		},
		{
			name: "invalid - missing firstName",
			user: &models.User{
				UUID:  "123e4567-e89b-12d3-a456-426614174000",
				Email: "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "invalid - bad email format",
			user: &models.User{
				UUID:      "123e4567-e89b-12d3-a456-426614174000",
				Email:     "invalid-email",
				FirstName: "John",
			},
			wantErr: true,
			errMsg:  "invalid email format",
		},
		{
			name: "invalid - invalid sex",
			user: &models.User{
				UUID:      "123e4567-e89b-12d3-a456-426614174000",
				Email:     "test@example.com",
				FirstName: "John",
				Sex:       "invalid",
			},
			wantErr: true,
			errMsg:  "sex must be one of: female, male, other",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateUserForRegistration(tt.user)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// Test_validateRequiredFields тестирует валидацию обязательных полей.
func Test_validateRequiredFields(t *testing.T) {
	tests := []struct {
		name    string
		user    *models.User
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid - all required fields present",
			user: &models.User{
				UUID:  "123e4567-e89b-12d3-a456-426614174000",
				Email: "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "invalid - empty UUID",
			user: &models.User{
				UUID:      "",
				Email:     "test@example.com",
				FirstName: "John",
			},
			wantErr: true,
			errMsg:  "user UUID is required",
		},
		{
			name: "invalid - whitespace only UUID",
			user: &models.User{
				UUID:      "   ",
				Email:     "test@example.com",
				FirstName: "John",
			},
			wantErr: true,
			errMsg:  "user UUID is required",
		},
		{
			name: "invalid - empty email",
			user: &models.User{
				UUID:      "123e4567-e89b-12d3-a456-426614174000",
				Email:     "",
				FirstName: "John",
			},
			wantErr: true,
			errMsg:  "email is required",
		},
		{
			name: "invalid - whitespace only email",
			user: &models.User{
				UUID:      "123e4567-e89b-12d3-a456-426614174000",
				Email:     "   ",
				FirstName: "John",
			},
			wantErr: true,
			errMsg:  "email is required",
		},
		{
			name: "invalid - malformed email",
			user: &models.User{
				UUID:      "123e4567-e89b-12d3-a456-426614174000",
				Email:     "not-an-email",
				FirstName: "John",
			},
			wantErr: true,
			errMsg:  "invalid email format",
		},
		{
			name: "invalid - email too long",
			user: &models.User{
				UUID:      "123e4567-e89b-12d3-a456-426614174000",
				Email:     strings.Repeat("a", 250) + "@example.com", // >255 chars
				FirstName: "John",
			},
			wantErr: true,
			errMsg:  "email is too long",
		},
		{
			name: "invalid - empty firstName",
			user: &models.User{
				UUID:      "123e4567-e89b-12d3-a456-426614174000",
				Email:     "test@example.com",
				FirstName: "",
			},
			wantErr: false,
		},
		{
			name: "invalid - whitespace only firstName",
			user: &models.User{
				UUID:      "123e4567-e89b-12d3-a456-426614174000",
				Email:     "test@example.com",
				FirstName: "   ",
			},
			wantErr: false,
		},
		{
			name: "valid - firstName boundary - 2 chars",
			user: &models.User{
				UUID:      "123e4567-e89b-12d3-a456-426614174000",
				Email:     "test@example.com",
				FirstName: "Jo",
			},
			wantErr: false,
		},
		{
			name: "valid - firstName boundary - 100 chars",
			user: &models.User{
				UUID:      "123e4567-e89b-12d3-a456-426614174000",
				Email:     "test@example.com",
				FirstName: strings.Repeat("A", 100),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateRequiredFields(tt.user)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// Test_validateOptionalFields тестирует валидацию опциональных полей.
func Test_validateOptionalFields(t *testing.T) {
	tests := []struct {
		name    string
		user    *models.User
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid - empty optional fields",
			user:    &models.User{},
			wantErr: false,
		},
		{
			name: "valid - all optional fields filled correctly",
			user: &models.User{
				LastName: "Doe",
				Sex:      "female",
				City:     "Moscow",
				Country:  "Russia",
			},
			wantErr: false,
		},
		{
			name: "valid - firstName filled correctly",
			user: &models.User{
				FirstName: "John",
			},
			wantErr: false,
		},
		{
			name: "invalid - firstName too short",
			user: &models.User{
				FirstName: "J",
			},
			wantErr: true,
			errMsg:  "firstName must be between 2 and 100 characters",
		},
		{
			name: "invalid - firstName too long",
			user: &models.User{
				FirstName: strings.Repeat("A", 101),
			},
			wantErr: true,
			errMsg:  "firstName must be between 2 and 100 characters",
		},
		{
			name: "valid - firstName boundary - 2 chars",
			user: &models.User{
				FirstName: "Jo",
			},
			wantErr: false,
		},
		{
			name: "valid - firstName boundary - 100 chars",
			user: &models.User{
				FirstName: strings.Repeat("A", 100),
			},
			wantErr: false,
		},
		{
			name: "invalid - lastName too short",
			user: &models.User{
				LastName: "D",
			},
			wantErr: true,
			errMsg:  "lastName must be between 2 and 100 characters",
		},
		{
			name: "invalid - lastName too long",
			user: &models.User{
				LastName: strings.Repeat("D", 101),
			},
			wantErr: true,
			errMsg:  "lastName must be between 2 and 100 characters",
		},
		{
			name: "valid - lastName boundary cases",
			user: &models.User{
				LastName: "Do", // 2 chars - minimum
			},
			wantErr: false,
		},
		{
			name: "valid - lastName 100 chars",
			user: &models.User{
				LastName: strings.Repeat("D", 100),
			},
			wantErr: false,
		},
		{
			name: "invalid - invalid sex value",
			user: &models.User{
				Sex: "invalid",
			},
			wantErr: true,
			errMsg:  "sex must be one of: female, male, other",
		},
		{
			name: "valid - sex female",
			user: &models.User{
				Sex: "female",
			},
			wantErr: false,
		},
		{
			name: "valid - sex male",
			user: &models.User{
				Sex: "male",
			},
			wantErr: false,
		},
		{
			name: "valid - sex other",
			user: &models.User{
				Sex: "other",
			},
			wantErr: false,
		},
		{
			name: "invalid - city too long",
			user: &models.User{
				City: strings.Repeat("M", 51),
			},
			wantErr: true,
			errMsg:  "city is too long",
		},
		{
			name: "valid - city 100 chars",
			user: &models.User{
				City: strings.Repeat("M", 50),
			},
			wantErr: false,
		},
		{
			name: "invalid - country too long",
			user: &models.User{
				Country: strings.Repeat("R", 51),
			},
			wantErr: true,
			errMsg:  "country is too long",
		},
		{
			name: "valid - country 100 chars",
			user: &models.User{
				Country: strings.Repeat("R", 50),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateOptionalFields(tt.user)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// Test_validateUserForRegistration_Integration интеграционный тест.
func Test_validateUserForRegistration_Integration(t *testing.T) {
	// Тест комбинированных ошибок - и обязательных, и опциональных полей
	user := &models.User{
		UUID:      "", // Обязательное поле отсутствует
		Email:     "test@example.com",
		FirstName: "John",
		Sex:       "invalid", // Невалидное опциональное поле
	}

	err := service.ValidateUserForRegistration(user)
	require.Error(t, err)
	// Должна быть ошибка обязательного поля (проверяется первым)
	assert.Contains(t, err.Error(), "user UUID is required")
}

// Бенчмарк тесты.
func Benchmark_validateUserForRegistration(b *testing.B) {
	user := &models.User{
		UUID:      "123e4567-e89b-12d3-a456-426614174000",
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Sex:       "male",
		City:      "Moscow",
		Country:   "Russia",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.ValidateUserForRegistration(user)
	}
}

// Пример теста с использованием table-driven подхода для разных email форматов.
func Test_validateRequiredFields_EmailFormats(t *testing.T) {
	baseUser := &models.User{
		UUID:      "123e4567-e89b-12d3-a456-426614174000",
		FirstName: "John",
	}

	emailTests := []struct {
		email   string
		valid   bool
		message string
	}{
		{"test@example.com", true, "standard email"},
		{"user.name@example.com", true, "email with dot"},
		{"user+tag@example.com", true, "email with plus"},
		{"test@example.co.uk", true, "email with multiple domains"},
		{"", false, "empty email"},
		{"invalid", false, "no @ symbol"},
		{"@example.com", false, "missing local part"},
		{"test@", false, "missing domain"},
		{"test @example.com", false, "space in email"},
		{"test..test@example.com", false, "double dot"},
	}

	for _, tt := range emailTests {
		t.Run(tt.message, func(t *testing.T) {
			user := *baseUser // Copy
			user.Email = tt.email

			err := service.ValidateRequiredFields(&user)
			if tt.valid {
				assert.NoError(t, err, "Expected %s to be valid", tt.email)
			} else {
				assert.Error(t, err, "Expected %s to be invalid", tt.email)
			}
		})
	}
}
