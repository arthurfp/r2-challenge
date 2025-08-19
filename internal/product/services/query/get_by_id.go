package query

import (
	"context"

	"r2-challenge/internal/product/domain"
)

type GetByIDService interface {
	GetByID(ctx context.Context, id string) (domain.Product, error)
}

func (s *service) GetByID(ctx context.Context, id string) (domain.Product, error) {
	ctx, span := s.tracer.StartSpan(ctx, "ProductQuery.GetByID")
	defer span.End()

	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return domain.Product{}, err
	}

	return p, nil
}
