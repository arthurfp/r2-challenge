package db

import (
    "context"

    "r2-challenge/internal/product/domain"
)

type ProductFilter struct {
    Category string
    Name     string
    Limit    int
    Offset   int
    SortBy   string
    SortDesc bool
}

type ProductRepository interface {
    Save(ctx context.Context, p domain.Product) (domain.Product, error)
    Update(ctx context.Context, p domain.Product) (domain.Product, error)
    Delete(ctx context.Context, id string) error
    GetByID(ctx context.Context, id string) (domain.Product, error)
    List(ctx context.Context, f ProductFilter) ([]domain.Product, error)
}


