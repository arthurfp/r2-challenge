package notification

import "context"

type Sender interface {
	SendOrderConfirmation(ctx context.Context, toEmail string, orderID string) error
}
