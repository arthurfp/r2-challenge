package query

import (
	repo "r2-challenge/internal/product/adapters/db"
	"r2-challenge/pkg/observability"
)

type service struct {
	repo   repo.ProductRepository
	tracer observability.Tracer
}

func NewService(r repo.ProductRepository, t observability.Tracer) (GetByIDService, ListService, error) {
	return &service{repo: r, tracer: t}, &service{repo: r, tracer: t}, nil
}
