package command

import (
	"context"

	repo "r2-challenge/internal/product/adapters/db"
	"r2-challenge/pkg/observability"
)

type DeleteService interface {
	Delete(ctx context.Context, productID string) error
}

type deleteService struct {
	repo   repo.ProductRepository
	tracer observability.Tracer
}

func NewDeleteService(r repo.ProductRepository, t observability.Tracer) (DeleteService, error) {
	return &deleteService{repo: r, tracer: t}, nil
}

func (s *deleteService) Delete(ctx context.Context, productID string) error {
	ctx, span := s.tracer.StartSpan(ctx, "ProductCommand.Delete")
	defer span.End()

	if err := s.repo.Delete(ctx, productID); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}
