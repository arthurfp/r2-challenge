package query

import (
    "context"

    repo "r2-challenge/internal/order/adapters/db"
    "r2-challenge/internal/order/domain"
    "r2-challenge/pkg/observability"
)

type GetByIDService interface { GetByID(ctx context.Context, id string) (domain.Order, error) }

type service struct {
    repo   repo.OrderRepository
    tracer observability.Tracer
}

func NewService(r repo.OrderRepository, t observability.Tracer) (GetByIDService, error) {
    return &service{repo: r, tracer: t}, nil
}

func (s *service) GetByID(ctx context.Context, id string) (domain.Order, error) {
    ctx, span := s.tracer.StartSpan(ctx, "OrderQuery.GetByID")
    defer span.End()
    
	order, err := s.repo.GetByID(ctx, id)
    if err != nil {
        span.RecordError(err)
        return domain.Order{}, err
    }
    
	return order, nil
}


