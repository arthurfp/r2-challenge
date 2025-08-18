# Deployment strategy (suggestion)

## Containerization
- Use the provided `Dockerfile` (multi-stage, distroless final image)
- Build locally: `make docker-build`
- Run locally: `make docker-run`

## Environment variables
- See `README.md` for all envs. Minimum for prod:
  - DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE
  - JWT_SECRET, JWT_ISSUER, JWT_EXPIRE
  - RATE_LIMIT_RPM, METRICS_* (optional), TLS_CERT_FILE/TLS_KEY_FILE (optional)

## Options
- Single container behind a reverse proxy (NGINX/Traefik) handling TLS
- Kubernetes (preferred for scale):
  - Deployment with 2+ replicas, HPA on CPU/RAM
  - Service (ClusterIP) + Ingress
  - ConfigMap for non-secrets, Secret for credentials
  - PodDisruptionBudget and Readiness/Liveness probes
  - Prometheus Operator scraping `/metrics`
  - RBAC: enforce admin-only routes at the API level (already implemented) and at the Ingress
  - Rolling updates and surge settings to ensure zero-downtime

## Postgres
- Managed DB (e.g., RDS/GCP Cloud SQL) with TLS and automated backups
- Tune max connections and indexes according to workload

## Observability
- OTEL tracing (already setup)
- Prometheus scraping `/metrics` + Grafana dashboards
- Structured logs (zap) to your log aggregator (ELK/Cloud Logging)

## Security
- JWT with short expiration and secret rotation
- Rate limiting enabled
- TLS enforced at the edge (proxy or app)

## CI/CD
- Pipeline: build, test, lint, image build, push, deploy
- Use immutable tags via `SERVICE`/`VERSION`
 - Run migrations as a separate job or initContainer

## What the challenge asked (addressed here)
- Containerization with Docker (Dockerfile, compose)
- Deployment suggestions (reverse proxy, Kubernetes, managed DB)
- Observability (OTEL tracing, Prometheus metrics)
- Security (JWT, RBAC, rate limiting, TLS)
