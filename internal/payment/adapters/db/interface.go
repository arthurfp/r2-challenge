package db

import (
    "context"
    pmtdomain "r2-challenge/internal/payment/domain"
)

type Repository interface {
    Save(ctx context.Context, payment pmtdomain.Payment) (pmtdomain.Payment, error)
}


