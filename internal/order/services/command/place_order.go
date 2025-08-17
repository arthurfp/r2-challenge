package command

import (
	"context"

	orderdb "r2-challenge/internal/order/adapters/db"
	"r2-challenge/internal/order/adapters/notification"
	"r2-challenge/internal/order/adapters/payment"
	"r2-challenge/internal/order/domain"
	pmtdomain "r2-challenge/internal/payment/domain"
	pmtcmd "r2-challenge/internal/payment/services/command"
	"r2-challenge/pkg/observability"
)

type PlaceOrderService interface { Place(ctx context.Context, order domain.Order) (domain.Order, error) }

type placeOrderService struct {
	repo        orderdb.OrderRepository
	payments    payment.Processor
	notifier    notification.Sender
	paymentsSvc pmtcmd.RecordService
	tracer      observability.Tracer
}

func NewPlaceOrderService(r orderdb.OrderRepository, p payment.Processor, n notification.Sender, t observability.Tracer, pr pmtcmd.RecordService) (PlaceOrderService, error) {
	return &placeOrderService{repo: r, payments: p, notifier: n, tracer: t, paymentsSvc: pr}, nil
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

	receiptID, err := s.payments.Charge(ctx, saved.UserID, saved.TotalCents)
	if err != nil {
		span.RecordError(err)
		return domain.Order{}, err
	}

	paymentRecord := pmtdomain.Payment{
		OrderID:     saved.ID,
		UserID:      saved.UserID,
		AmountCents: saved.TotalCents,
		Provider:    "mock",
		ReceiptID:   receiptID,
		Status:      "captured",
	}

	_, _ = s.paymentsSvc.Record(ctx, paymentRecord)

	_ = s.notifier.SendOrderConfirmation(ctx, "user@example.com", saved.ID)

	return saved, nil
}


