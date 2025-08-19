package db

import (
	"context"

	"r2-challenge/internal/user/domain"
)

type UserFilter struct {
	Email  string
	Name   string
	Limit  int
	Offset int
}

type UserRepository interface {
	Save(ctx context.Context, u domain.User) (domain.User, error)
	Update(ctx context.Context, u domain.User) (domain.User, error)
	GetByID(ctx context.Context, id string) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	List(ctx context.Context, f UserFilter) ([]domain.User, error)
}
