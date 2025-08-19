package command

import (
	"context"

	repo "r2-challenge/internal/order/adapters/db"
	"r2-challenge/internal/order/domain"
	"r2-challenge/pkg/observability"
)

type UpdateStatusService interface {
	UpdateStatus(ctx context.Context, orderID string, status string) (domain.Order, error)
}

type updateStatusService struct {
	repo   repo.OrderRepository
	tracer observability.Tracer
}

func NewUpdateStatusService(r repo.OrderRepository, t observability.Tracer) (UpdateStatusService, error) {
	return &updateStatusService{repo: r, tracer: t}, nil
}

func (s *updateStatusService) UpdateStatus(ctx context.Context, orderID string, status string) (domain.Order, error) {
	ctx, span := s.tracer.StartSpan(ctx, "OrderCommand.UpdateStatus")
	defer span.End()

	order, err := s.repo.UpdateStatus(ctx, orderID, status)
	if err != nil {
		span.RecordError(err)
		return domain.Order{}, err
	}

	return order, nil
}
