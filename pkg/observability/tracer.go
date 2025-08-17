package observability

import (
	"context"
)

type Tracer interface {
	StartSpan(ctx context.Context, name string) (context.Context, func())
}

type noopTracer struct{}

func (noopTracer) StartSpan(ctx context.Context, _ string) (context.Context, func()) {
	return ctx, func() {}
}

// SetupTracer returns a no-op tracer; replace with OTEL later if needed.
func SetupTracer() (Tracer, error) {
	return noopTracer{}, nil
}
