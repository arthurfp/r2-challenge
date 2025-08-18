package payment

import "context"

type noopProcessor struct{}

func NewNoopProcessor() Processor { return noopProcessor{} }

func (noopProcessor) Charge(_ context.Context, _ string, _ int64) (string, error) {
    return "noop-receipt", nil
}


