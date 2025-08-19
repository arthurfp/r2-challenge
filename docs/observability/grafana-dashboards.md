## Grafana Dashboards (example)

- Prometheus scrape endpoint: `/metrics` (or dedicated port if configured).
- Suggested dashboards:
  - Go runtime metrics (Go / Prometheus)
  - Echo HTTP metrics: `http_server_duration_seconds`, `http_requests_total`
  - OTEL traces to visualize latency outliers

### Example Prometheus queries

- P95 latency by route:
```
histogram_quantile(0.95, sum(rate(http_server_request_duration_seconds_bucket[5m])) by (le, handler))
```

- Error rate by route:
```
sum(rate(http_requests_total{code!~"2.."}[5m])) by (handler) / sum(rate(http_requests_total[5m])) by (handler)
```


