package domain

import "time"

type Product struct {
    ID          string    `json:"id"`
    Name        string    `json:"name" validate:"required,min=3"`
    Description string    `json:"description"`
    Category    string    `json:"category" validate:"required"`
    PriceCents  int64     `json:"price_cents" validate:"required,gte=0"`
    Inventory   int64     `json:"inventory" validate:"gte=0"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    DeletedAt  *time.Time  
}


