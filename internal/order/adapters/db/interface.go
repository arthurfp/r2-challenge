package db

import (
	"context"

	"r2-challenge/internal/order/domain"
)

type OrderFilter struct {
	UserID string
	Limit  int
	Offset int
}

type OrderRepository interface {
	Save(ctx context.Context, order domain.Order) (domain.Order, error)
	UpdateStatus(ctx context.Context, orderID string, status string) (domain.Order, error)
	GetByID(ctx context.Context, orderID string) (domain.Order, error)
	ListByUser(ctx context.Context, userID string, filter OrderFilter) ([]domain.Order, error)
}
