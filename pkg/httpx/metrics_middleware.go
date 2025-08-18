package httpx

import (
    "time"

    "github.com/labstack/echo/v4"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/metric"
)

// MetricsMiddleware instruments HTTP requests with OTEL metrics.
func MetricsMiddleware(meter metric.Meter) echo.MiddlewareFunc {
    requestCounter, _ := meter.Int64Counter("http_requests_total")
    inFlight, _ := meter.Int64UpDownCounter("http_in_flight_requests")
    durationHist, _ := meter.Float64Histogram("http_request_duration_seconds")

    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            ctx := c.Request().Context()
            method := c.Request().Method
            route := c.Path()

            inFlight.Add(ctx, 1)
            defer inFlight.Add(ctx, -1)

            start := time.Now()
            err := next(c)
            elapsed := time.Since(start).Seconds()

            status := c.Response().Status
            attrs := []attribute.KeyValue{
                attribute.String("http.method", method),
                attribute.String("http.route", route),
                attribute.Int("http.status_code", status),
            }
            requestCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
            durationHist.Record(ctx, elapsed, metric.WithAttributes(attrs...))

            return err
        }
    }
}


