package domain

import "time"

type User struct {
	ID           string     `json:"id"`
	Email        string     `json:"email" validate:"required,email"`
	Name         string     `json:"name" validate:"required,min=3"`
	PasswordHash string     `json:"-"`
	Role         string     `json:"role" validate:"required"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}
