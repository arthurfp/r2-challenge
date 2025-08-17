package payment

import (
    "context"
    "fmt"
    "time"

    "r2-challenge/pkg/observability"
)

type mockProcessor struct { tracer observability.Tracer }

func NewMockProcessor(t observability.Tracer) (Processor, error) {
    return &mockProcessor{tracer: t}, nil
}

func (m *mockProcessor) Charge(ctx context.Context, userID string, amountCents int64) (string, error) {
    ctx, span := m.tracer.StartSpan(ctx, "Payment.Charge")
    defer span.End()

    // Simulate external processing latency
    time.Sleep(50 * time.Millisecond)

    // Always succeed for now
    receipt := fmt.Sprintf("rcpt_%d", time.Now().UnixNano())
    return receipt, nil
}


