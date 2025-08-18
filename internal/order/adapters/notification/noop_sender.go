package notification

import "context"

type noopSender struct{}

func NewNoopSender() Sender { return noopSender{} }

func (noopSender) SendOrderConfirmation(_ context.Context, _ string, _ string) error { return nil }


