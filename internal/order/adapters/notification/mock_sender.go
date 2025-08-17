package notification

import (
    "context"
    "time"

    "r2-challenge/pkg/observability"
)

type mockSender struct { tracer observability.Tracer }

func NewMockSender(t observability.Tracer) (Sender, error) {
    return &mockSender{tracer: t}, nil
}

func (m *mockSender) SendOrderConfirmation(ctx context.Context, toEmail string, orderID string) error {
    ctx, span := m.tracer.StartSpan(ctx, "Notification.SendOrderConfirmation")
    defer span.End()

    _ = toEmail
    _ = orderID

    time.Sleep(10 * time.Millisecond)
    return nil
}


