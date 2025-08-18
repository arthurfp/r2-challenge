package metrics

import (
    "net/http"

    prom "github.com/prometheus/client_golang/prometheus"
    promhttp "github.com/prometheus/client_golang/prometheus/promhttp"
    "go.opentelemetry.io/otel"
    otelprom "go.opentelemetry.io/otel/exporters/prometheus"
    sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// Setup creates an OTEL Meter and a Prometheus exporter handler.
func Setup() (otelMeterProvider *sdkmetric.MeterProvider, handler http.Handler, err error) {
    reg := prom.NewRegistry()

    exp, err := otelprom.New(otelprom.WithRegisterer(reg))
    if err != nil {
        return nil, nil, err
    }

    provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(exp))
    otel.SetMeterProvider(provider)

    return provider, promhttp.HandlerFor(reg, promhttp.HandlerOpts{}), nil
}


