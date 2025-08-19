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
 - Caching: Redis decorators for read-heavy queries (Products/Users) with on-write invalidation
 - Idempotency: middleware for POST /orders (deduplicates retries with Idempotency-Key)
 - Concurrency: atomic stock decrement during order placement
 - Performance tooling: k6 script for smoke/load and example Grafana queries

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
```

## Environment (essentials)
- HTTP: `HTTP_HOST`, `HTTP_PORT`, `READ_HEADER_TIMEOUT`, `HTTP_TIMEOUT`
- DB: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`
- JWT: `JWT_SECRET`, `JWT_ISSUER`, `JWT_EXPIRE`
- Rate limit: `RATE_LIMIT_RPM`
- Metrics: `METRICS_ENABLED`, `METRICS_PATH`, `METRICS_PORT`
- TLS (optional): `TLS_CERT_FILE`, `TLS_KEY_FILE`
 - Redis (optional, enables caching + idempotency storage): `REDIS_ADDR` (e.g. `localhost:6379`), `REDIS_PASSWORD`, `REDIS_DB` (default `0`)

go run ./cmd/app

## OpenAPI / Swagger
- Spec generated from annotations in HTTP handlers: `make swaggen` (alias of `swag-gen-app`)
- Output file: `cmd/app/swagger-gen/swagger.yaml`
- Live docs: once running, open `http://localhost:8080/swagger` and try endpoints directly in the browser

## Caching (Redis)
Caching is opt-in. If `REDIS_ADDR` is provided, Product and User repositories are wrapped with a cache decorator:

- Keys: simple composite strings by id or filter (scope per module)
- TTL: short and conservative by default
- Invalidation: on writes (create/update/delete), related keys are deleted; lists use a namespaced prefix

Recommended to start with Redis locally and short TTLs; expand as access patterns stabilize.

## Metrics (OTEL/Prometheus)
- Enabled by default at `GET /metrics` (or dedicated port with `METRICS_PORT`)
- Scrape example:
```yaml
scrape_configs:
  - job_name: 'r2-challenge'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: /metrics
```

## Deployment (short)
- Dockerfile (multi-stage, distroless) + compose
- Suggested K8s setup in `docs/deployment.md` (replicas, HPA, Ingress, managed Postgres, Prometheus)

## Idempotency & concurrency

- Idempotency (POST /orders): clients send an `Idempotency-Key` header and may retry the same payload. The server stores status/response in Redis for a short TTL and returns the exact same response for the same key+payload. Different payload under the same key responds with `409 Conflict`.

- Atomic inventory decrement: order placement includes a conditional update (`inventory = inventory - ? WHERE id = ? AND inventory >= ?`), guaranteeing no negative stock under concurrent requests. If the condition fails, the transaction is rolled back.

- Locking strategy: optimistic-first (timestamps/versions) is a natural next step if product updates become highly concurrent; pessimistic locking is reserved for hotspots.

## Load testing (k6)

A minimal script is included to exercise public reads and order placement with idempotency:

```bash
# tokens are optional (only needed for private endpoints)
K6_BASE_URL=http://localhost:8080 \
K6_USER_TOKEN=<jwt> \
K6_ADMIN_TOKEN=<jwt> \
k6 run scripts/k6/order_place.js
```

The script also retries an order with the same Idempotency-Key to validate deduplication.

## Dashboards (Grafana)

See `docs/observability/grafana-dashboards.md` for example Prometheus queries (p95 latency per route, error rate) and suggested dashboards.

## Notes on style and decisions
- Single domain type per module; services operate only over domain types
- Handlers are thin; all routes registered in `main`
- Early returns, explicit error contexts, descriptive names
- Timestamps handled in DB adapter only (no duplication in services)

## Where to read more
- API docs: `docs/api/products.md`, `docs/api/users.md`, `docs/api/orders.md`
- Deployment: `docs/deployment.md`
