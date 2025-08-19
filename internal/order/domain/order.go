package domain

import "time"

// Order is the core domain type for orders.
type Order struct {
	ID         string      `json:"id"`
	UserID     string      `json:"user_id" validate:"required"`
	Status     string      `json:"status"`
	TotalCents int64       `json:"total_cents"`
	Items      []OrderItem `json:"items"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	DeletedAt  *time.Time  `json:"deleted_at"`
}

type OrderItem struct {
	ID         string     `json:"id"`
	OrderID    string     `json:"order_id"`
	ProductID  string     `json:"product_id" validate:"required"`
	Quantity   int64      `json:"quantity" validate:"required,gt=0"`
	PriceCents int64      `json:"price_cents" validate:"required,gte=0"`
	DeletedAt  *time.Time `json:"deleted_at"`
}
