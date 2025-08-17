package observability

import (
	"context"
	"io"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type Span interface {
	End()
	RecordError(err error)
}

type Tracer interface {
	StartSpan(ctx context.Context, name string) (context.Context, Span)
}

type otelTracer struct{ tracer trace.Tracer }

type otelSpan struct{ span trace.Span }

func (s otelSpan) End() { s.span.End() }
func (s otelSpan) RecordError(err error) {
	if err == nil { return }
	s.span.RecordError(err)
	s.span.SetAttributes(attribute.String("error", err.Error()))
}

func (t otelTracer) StartSpan(ctx context.Context, name string) (context.Context, Span) {
	ctx, sp := t.tracer.Start(ctx, name)
	return ctx, otelSpan{span: sp}
}

// SetupTracer sets a basic OTEL tracer provider; stdout exporter is disabled by default.
func SetupTracer() (Tracer, error) {
	exp, err := stdouttrace.New(stdouttrace.WithWriter(io.Discard))
	if err != nil { return nil, err }
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exp))
	otel.SetTracerProvider(tp)
	return otelTracer{tracer: otel.Tracer("r2-challenge")}, nil
}
