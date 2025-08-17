package command

import (
    "context"

    repo "r2-challenge/internal/product/adapters/db"
    "r2-challenge/internal/product/domain"
    "r2-challenge/pkg/observability"
)

type CreateService interface {
    Create(ctx context.Context, p domain.Product) (domain.Product, error)
}

type createService struct {
    repo   repo.ProductRepository
    tracer observability.Tracer
}

func NewCreateService(r repo.ProductRepository, t observability.Tracer) (CreateService, error) {
    return &createService{repo: r, tracer: t}, nil
}

func (s *createService) Create(ctx context.Context, p domain.Product) (domain.Product, error) {
    ctx, span := s.tracer.StartSpan(ctx, "ProductCommand.Create")
    defer span.End()

    res, err := s.repo.Save(ctx, p)
    if err != nil {
        span.RecordError(err)
        return res, err
    }

    return res, nil
}


