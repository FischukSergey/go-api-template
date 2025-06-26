package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/Fisher-Development/woman-app-backend/internal/models"
)

var (
	ErrUserAlreadyExists  = errors.New("user with this UUID already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

// CreateUser creates a new user.
func (s *Storage) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, email, first_name)
		VALUES ($1, $2, $3)
	`
	_, err := s.db.Exec(ctx, query, user.UUID, user.Email, user.FirstName)
	if err != nil {
		// Проверяем на ошибку дублирования UUID/email
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				if strings.Contains(pgErr.Detail, "id") {
					return ErrUserAlreadyExists
				}
				if strings.Contains(pgErr.Detail, "email") {
					return ErrEmailAlreadyExists
				}
			}
		}
		return err
	}
	return nil
}

// GetUserByUUID returns a user by UUID.
func (s *Storage) GetUserByUUID(ctx context.Context, uuid string) (*models.User, error) {
	query := `
		SELECT 
			id, 
			email, 
			first_name, 
			last_name, 
			sex, 
			city, 
			country, 
			date_of_birth, 
			created_at, 
			updated_at
		FROM users
		WHERE id = $1
	`
	row := s.db.QueryRow(ctx, query, uuid)
	var user models.User
	var birthDateStr, lastNameStr, sexStr, cityStr, countryStr sql.NullString

	err := row.Scan(
		&user.UUID,
		&user.Email,
		&user.FirstName,
		&lastNameStr,
		&sexStr,
		&cityStr,
		&countryStr,
		&birthDateStr,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Присваиваем значения из NullString
	if birthDateStr.Valid {
		user.BirthDate = birthDateStr.String
	}
	if lastNameStr.Valid {
		user.LastName = lastNameStr.String
	}
	if sexStr.Valid {
		user.Sex = sexStr.String
	}
	if cityStr.Valid {
		user.City = cityStr.String
	}
	if countryStr.Valid {
		user.Country = countryStr.String
	}

	return &user, nil
}

// UpdateUser updates a user profile.
func (s *Storage) UpdateUser(ctx context.Context, user *models.User) error {
	//проверяем есть ли пользователь в базе данных
	_, err := s.GetUserByUUID(ctx, user.UUID)
	if err != nil {
		return err
	}

	// Подготавливаем birthDate для БД
	var birthDate any
	if user.BirthDate == "" {
		birthDate = nil
	} else {
		birthDate = user.BirthDate
	}

	//обновляем пользователя
	query := `
		UPDATE users
		SET 
			first_name = $1, 
			last_name = $2, 
			sex = $3, 
			city = $4, 
			country = $5, 
			date_of_birth = $6
		WHERE id = $7
	`
	_, err = s.db.Exec(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.Sex,
		user.City,
		user.Country,
		birthDate,
		user.UUID,
	)
	return err
}
