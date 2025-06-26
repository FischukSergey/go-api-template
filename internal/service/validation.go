package service

import (
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/Fisher-Development/woman-app-backend/internal/models"
)

// ValidateUserForRegistration валидирует данные пользователя для регистрации.
func ValidateUserForRegistration(user *models.User) error {
	if err := ValidateRequiredFields(user); err != nil {
		return err
	}
	return ValidateOptionalFields(user)
}

// ValidateUserForUpdate валидирует данные пользователя для обновления.
func ValidateUserForUpdate(user *models.User) error {
	return ValidateOptionalFields(user)
}

// ValidateRequiredFields проверяет обязательные поля.
func ValidateRequiredFields(user *models.User) error {
	// UUID должен быть установлен
	if strings.TrimSpace(user.UUID) == "" {
		return errors.New("user UUID is required")
	}

	// Email обязателен
	if strings.TrimSpace(user.Email) == "" {
		return errors.New("email is required")
	}
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return errors.New("invalid email format")
	}
	if len(user.Email) > 255 {
		return errors.New("email is too long (max 255 characters)")
	}

	// Имя обязательно
	// if strings.TrimSpace(user.FirstName) == "" {
	// 	return errors.New("firstName is required")
	// }
	// if len(user.FirstName) < 2 || len(user.FirstName) > 100 {
	// 	return errors.New("firstName must be between 2 and 100 characters")
	// }
	return nil
}

// validateOptionalFields проверяет опциональные поля.
func ValidateOptionalFields(user *models.User) error {
	// Имя (опционально)
	if user.FirstName != "" && (len(user.FirstName) < 2 || len(user.FirstName) > 100) {
		return errors.New("firstName must be between 2 and 100 characters")
	}

	// Фамилия (опциональная)
	if user.LastName != "" && (len(user.LastName) < 2 || len(user.LastName) > 100) {
		return errors.New("lastName must be between 2 and 100 characters")
	}

	// Пол (опциональный)
	if user.Sex != "" {
		validSex := map[string]bool{"female": true, "male": true, "other": true}
		if !validSex[user.Sex] {
			return errors.New("sex must be one of: female, male, other")
		}
	}

	// Город (опциональный)
	if user.City != "" && len(user.City) > 50 {
		return errors.New("city is too long (max 50 characters)")
	}

	// Страна (опциональная)
	if user.Country != "" && len(user.Country) > 50 {
		return errors.New("country is too long (max 50 characters)")
	}

	// Валидация birthDate (опционально)
	if user.BirthDate != "" {
		if _, err := time.Parse("2006-01-02", user.BirthDate); err != nil {
			return errors.New("invalid birthDate format, expected YYYY-MM-DD")
		}
	}

	return nil
}
