# r2-challenge

Didactic e-commerce API using a hexagonal architecture, CQRS-style services, JWT/RBAC, GORM/Postgres, OTEL tracing, and tests with gomock.

## Run locally
```bash
make docker-up
make db-setup
make test
# to run the app
# go run ./cmd/app
 # or with Docker
 make docker-build && make docker-run

## Environment
- `HTTP_HOST` (default: localhost)
- `HTTP_PORT` (default: 8080)
- `READ_HEADER_TIMEOUT` (default: 15s)
- `HTTP_TIMEOUT` (default: 10s)
- `RATE_LIMIT_RPM` (default: 60)
- `TLS_CERT_FILE` (optional): Path to PEM certificate
- `TLS_KEY_FILE` (optional): Path to PEM private key
- `DB_*` (host, port, user, password, name, sslmode, pool params)
- `JWT_SECRET`, `JWT_ISSUER`, `JWT_EXPIRE`

### Example: run with TLS
Assuming you have `server.crt` and `server.key` in the project root:
```bash
# start the server
go run ./cmd/app

## Metrics (OTEL/Prometheus)
- `METRICS_ENABLED` (default: true)
- `METRICS_PATH` (default: /metrics)
- `METRICS_PORT` (optional): if set, metrics are exposed on a dedicated server at `:${METRICS_PORT}`

### Example Prometheus scrape config
Scraping from the main server endpoint:
```yaml
scrape_configs:
  - job_name: 'r2-challenge'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: /metrics
Or scraping from a dedicated metrics server:
```yaml
scrape_configs:
  - job_name: 'r2-challenge-metrics'
    static_configs:
      - targets: ['localhost:9090']

## Documentation
- Overview and guides: `docs/`
- Products: `docs/api/products.md`
- Users: `docs/api/users.md`
- Orders: `docs/api/orders.md`
- Deployment: `docs/deployment.md`

### OpenAPI / Swagger
- OpenAPI spec is generated from annotations on HTTP handlers.
- Generate/update spec with: `make swaggen` (alias to `swag-gen-app`).
- Output file: `cmd/app/swagger-gen/swagger.yaml`.

## Structure
- `cmd/app`: DI (fx), HTTP server and route registration
- `internal/<bounded_context>`: domain, adapters (db/http), services (command/query)
- `pkg`: infrastructure/utilities (auth, httpx, observability, validator, logger)
- `db/migrations`: SQL migrations

## WIP
