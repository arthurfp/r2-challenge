package command

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	orderdb "r2-challenge/internal/order/adapters/db"
	notifmock "r2-challenge/internal/order/adapters/notification"
	paymentmock "r2-challenge/internal/order/adapters/payment"
	"r2-challenge/internal/order/domain"
	pmtdomain "r2-challenge/internal/payment/domain"
	"r2-challenge/pkg/observability"
)

// stubRecordSvc matches the RecordService signature expected by PlaceOrderService.
type stubRecordSvc struct{}

func (stubRecordSvc) Record(_ context.Context, _ pmtdomain.Payment) (pmtdomain.Payment, error) {
	return pmtdomain.Payment{}, nil
}

func TestPlaceOrder_Success(t *testing.T) {
	tracer, _ := observability.SetupTracer()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := orderdb.NewMockOrderRepository(ctrl)
	payments := paymentmock.NewMockProcessor(ctrl)
	notifier := notifmock.NewMockSender(ctrl)

	s, err := NewPlaceOrderService(repo, payments, notifier, tracer, stubRecordSvc{})
	if err != nil {
		t.Fatalf("failed to build service: %v", err)
	}

	order := domain.Order{UserID: "u1", Items: []domain.OrderItem{{ProductID: "p1", Quantity: 1, PriceCents: 1000}}, TotalCents: 1000}

	repo.EXPECT().Save(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, o domain.Order) (domain.Order, error) { o.ID = "ord_1"; return o, nil })
	payments.EXPECT().Charge(gomock.Any(), "u1", int64(1000)).Return("rcpt_x", nil)
	notifier.EXPECT().SendOrderConfirmation(gomock.Any(), gomock.Any(), "ord_1").Return(nil)

	result, err := s.Place(context.Background(), order)
	if err != nil {
		t.Fatalf("Place failed: %v", err)
	}
	if result.ID == "" {
		t.Fatalf("expected order ID")
	}
}
