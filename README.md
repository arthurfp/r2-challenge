# r2-challenge

Didactic e-commerce API using a hexagonal architecture, CQRS-style services, JWT/RBAC, GORM/Postgres, OTEL tracing, and tests with gomock.

## Run locally
```bash
make docker-up
make db-setup
make test
# to run the app
# go run ./cmd/app

## Documentation
- Overview and guides: `docs/`
- Products: `docs/api/products.md`
- Users: `docs/api/users.md`
- Orders: `docs/api/orders.md`

## Structure
- `cmd/app`: DI (fx), HTTP server and route registration
- `internal/<bounded_context>`: domain, adapters (db/http), services (command/query)
- `pkg`: infrastructure/utilities (auth, httpx, observability, validator, logger)
- `db/migrations`: SQL migrations


## WIP
