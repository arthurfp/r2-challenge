package query

import (
	"context"

	repo "r2-challenge/internal/order/adapters/db"
	"r2-challenge/internal/order/domain"
)

type ListByUserService interface {
	ListByUser(ctx context.Context, userID string, filter repo.OrderFilter) ([]domain.Order, error)
}

func (s *service) ListByUser(ctx context.Context, userID string, filter repo.OrderFilter) ([]domain.Order, error) {
	ctx, span := s.tracer.StartSpan(ctx, "OrderQuery.ListByUser")
	defer span.End()

	list, err := s.repo.ListByUser(ctx, userID, filter)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return list, nil
}
