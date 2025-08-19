package payment

import "context"

type Processor interface {
	Charge(ctx context.Context, userID string, amountCents int64) (receiptID string, err error)
}
