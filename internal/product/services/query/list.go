package query

import (
    "context"

    repo "r2-challenge/internal/product/adapters/db"
    "r2-challenge/internal/product/domain"
)

type ListService interface {
    List(ctx context.Context, f repo.ProductFilter) ([]domain.Product, error)
}

func (s *service) List(ctx context.Context, f repo.ProductFilter) ([]domain.Product, error) {
    ctx, span := s.tracer.StartSpan(ctx, "ProductQuery.List")
    defer span.End()

    list, err := s.repo.List(ctx, f)
    if err != nil {
        span.RecordError(err)
        return nil, err
    }

    return list, nil
}


