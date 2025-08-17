package command

import (
    "context"

    repo "r2-challenge/internal/order/adapters/db"
    "r2-challenge/internal/order/domain"
    "r2-challenge/pkg/observability"
)

type PlaceOrderService interface { Place(ctx context.Context, order domain.Order) (domain.Order, error) }

type placeOrderService struct {
    repo   repo.OrderRepository
    tracer observability.Tracer
}

func NewPlaceOrderService(r repo.OrderRepository, t observability.Tracer) (PlaceOrderService, error) {
    return &placeOrderService{repo: r, tracer: t}, nil
}

func (s *placeOrderService) Place(ctx context.Context, order domain.Order) (domain.Order, error) {
    ctx, span := s.tracer.StartSpan(ctx, "OrderCommand.Place")
    defer span.End()

    // default status when not provided
    if order.Status == "" {
        order.Status = "created"
    }

    saved, err := s.repo.Save(ctx, order)
    if err != nil {
        span.RecordError(err)
        return domain.Order{}, err
    }

    return saved, nil
}


