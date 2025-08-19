package db

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
	pmtdomain "r2-challenge/internal/payment/domain"
	appdb "r2-challenge/pkg/db"
	"r2-challenge/pkg/observability"
)

type dbPaymentRepository struct {
	db     *gorm.DB
	tracer observability.Tracer
}

func NewDBRepository(database *appdb.Database, t observability.Tracer) (Repository, error) {
	return &dbPaymentRepository{db: database.DB, tracer: t}, nil
}

func (r *dbPaymentRepository) Save(ctx context.Context, payment pmtdomain.Payment) (pmtdomain.Payment, error) {
	ctx, span := r.tracer.StartSpan(ctx, "PaymentRepository.Save")
	defer span.End()

	now := time.Now().UTC()
	if payment.ID == "" {
		payment.ID = uuid.NewString()
	}
	payment.CreatedAt = now
	payment.UpdatedAt = now

	if err := r.db.WithContext(ctx).Table("payments").Create(&payment).Error; err != nil {
		span.RecordError(err)
		return pmtdomain.Payment{}, err
	}

	return payment, nil
}
