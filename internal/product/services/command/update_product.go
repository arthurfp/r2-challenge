package command

import (
    "context"

    repo "r2-challenge/internal/product/adapters/db"
    "r2-challenge/internal/product/domain"
    "r2-challenge/pkg/observability"
)

type UpdateService interface {
    Update(ctx context.Context, product domain.Product) (domain.Product, error)
}

type updateService struct {
    repo   repo.ProductRepository
    tracer observability.Tracer
}

func NewUpdateService(r repo.ProductRepository, t observability.Tracer) (UpdateService, error) {
    return &updateService{repo: r, tracer: t}, nil
}

func (s *updateService) Update(ctx context.Context, product domain.Product) (domain.Product, error) {
    ctx, span := s.tracer.StartSpan(ctx, "ProductCommand.Update")
    defer span.End()

    res, err := s.repo.Update(ctx, product)
    if err != nil {
        span.RecordError(err)
        return res, err
    }

    return res, nil
}


