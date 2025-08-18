# r2-challenge

A didactic e-commerce API built with a pragmatic hexagonal architecture. The goal was to keep it simple (KISS) while showcasing clean code, clear boundaries, and production-friendly concerns (auth, tracing, metrics, tests, deploy).

## What’s inside
- Hexagonal-ish layout: domain-first, adapters (db/http), and services split by intent (command/query)
- HTTP: Echo (thin handlers, routes in main)
- DI: Uber Fx
- DB: Postgres + GORM + SQL migrations
- AuthZ/AuthN: JWT + RBAC middleware
- Tracing: OpenTelemetry spans (errors recorded) + request span middleware
- Metrics: Prometheus endpoint (requests, duration, inflight)
- Rate limiting: per-IP token bucket
- Docs: OpenAPI (annotations) + live Swagger UI
- Tests: gomock-based unit tests (services) and build-ready project layout

## Folder structure (high-level)
- `cmd/app`: app composition (Fx), HTTP server and route registration only
- `internal/<bounded_context>`
  - `domain`: single domain type per module (shared across layers)
  - `adapters/db`: persistence interfaces and GORM repositories
  - `adapters/http`: one handler file per operation; validation plus mapping transport ↔ domain
  - `services/command|query`: one file per operation; no transport/infra deps
- `pkg`: infra/utilities (auth, observability, http helpers, validator, logger)
- `db/migrations`: SQL scripts

## Run locally
```bash
make docker-up      # starts Postgres (and app if you want via compose)
make db-setup       # waits for DB and applies migrations
make test           # unit tests
# run the app (envs default to local)
go run ./cmd/app
# or with Docker
make docker-build && make docker-run

## Environment (essentials)
- HTTP: `HTTP_HOST`, `HTTP_PORT`, `READ_HEADER_TIMEOUT`, `HTTP_TIMEOUT`
- DB: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`
- JWT: `JWT_SECRET`, `JWT_ISSUER`, `JWT_EXPIRE`
- Rate limit: `RATE_LIMIT_RPM`
- Metrics: `METRICS_ENABLED`, `METRICS_PATH`, `METRICS_PORT`
- TLS (optional): `TLS_CERT_FILE`, `TLS_KEY_FILE`

go run ./cmd/app

## OpenAPI / Swagger
- Spec generated from annotations in HTTP handlers: `make swaggen` (alias of `swag-gen-app`)
- Output file: `cmd/app/swagger-gen/swagger.yaml`
- Live docs: once running, open `http://localhost:8080/swagger` and try endpoints directly in the browser

## Metrics (OTEL/Prometheus)
- Enabled by default at `GET /metrics` (or dedicated port with `METRICS_PORT`)
- Scrape example:
```yaml
scrape_configs:
  - job_name: 'r2-challenge'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: /metrics

## Deployment (short)
- Dockerfile (multi-stage, distroless) + compose
- Suggested K8s setup in `docs/deployment.md` (replicas, HPA, Ingress, managed Postgres, Prometheus)

## Notes on style and decisions
- Single domain type per module; services operate only over domain types
- Handlers are thin; all routes registered in `main`
- Early returns, explicit error contexts, descriptive names
- Timestamps handled in DB adapter only (no duplication in services)

## Where to read more
- API docs: `docs/api/products.md`, `docs/api/users.md`, `docs/api/orders.md`
- Deployment: `docs/deployment.md`
