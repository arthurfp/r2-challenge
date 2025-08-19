package domain

import "time"

type Payment struct {
	ID          string    `json:"id"`
	OrderID     string    `json:"order_id"`
	UserID      string    `json:"user_id"`
	AmountCents int64     `json:"amount_cents"`
	Provider    string    `json:"provider"`
	ReceiptID   string    `json:"receipt_id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
