package command

import (
	"context"

	pmtdb "r2-challenge/internal/payment/adapters/db"
	pmtdomain "r2-challenge/internal/payment/domain"
	"r2-challenge/pkg/observability"
)

type RecordService interface {
	Record(ctx context.Context, payment pmtdomain.Payment) (pmtdomain.Payment, error)
}

type service struct {
	repo   pmtdb.Repository
	tracer observability.Tracer
}

func NewService(r pmtdb.Repository, t observability.Tracer) (RecordService, error) {
	return &service{repo: r, tracer: t}, nil
}

func (s *service) Record(ctx context.Context, payment pmtdomain.Payment) (pmtdomain.Payment, error) {
	ctx, span := s.tracer.StartSpan(ctx, "PaymentCommand.Record")
	defer span.End()

	saved, err := s.repo.Save(ctx, payment)
	if err != nil {
		span.RecordError(err)
		return pmtdomain.Payment{}, err
	}

	return saved, nil
}
