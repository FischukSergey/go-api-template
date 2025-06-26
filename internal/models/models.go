package models

import "time"

// User is a model for a user.
type User struct {
	UUID      string    `json:"uuid"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Phone     string    `json:"phone"`
	BirthDate string    `json:"birthDate"`
	Sex       string    `json:"sex"`
	City      string    `json:"city"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
