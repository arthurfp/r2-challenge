install-deps:
	# Installing Make Dependencies
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8
	go install github.com/swaggo/swag/cmd/swag@latest
	npm install @openapitools/openapi-generator-cli -g

run:
	go run ./cmd/app

build:
	go build ./...

tidy:
	go mod tidy

test:
	go test ./...

lint:
	# Running lint on project
	go mod tidy
	gofumpt -w .
	golangci-lint run --issues-exit-code 1 --out-format colored-tab

## Kill any process using TCP port 5432 (PostgreSQL default)
.PHONY: free-5432
free-5432:
	@echo "Releasing TCP port 5432 if occupied..."
	# Try stopping compose DB service first
	-@bash -c 'docker compose stop db 2>/dev/null || true'
	# Stop any container exposing host port 5432
	-@bash -c "docker ps --filter 'publish=5432' -q | xargs -r docker stop 2>/dev/null || true"
	# Kill local processes bound to 5432 (best-effort, may require privileges)
	-@bash -c 'fuser -k 5432/tcp 2>/dev/null || true'
	-@bash -c 'pids=$$(lsof -ti tcp:5432 2>/dev/null || true); [ -z "$$pids" ] || kill -9 $$pids 2>/dev/null || true'
	# Try stopping system service (passwordless sudo only)
	-@bash -c 'command -v systemctl >/dev/null && sudo -n systemctl stop postgresql 2>/dev/null || true'
	-@bash -c 'command -v service >/dev/null && sudo -n service postgresql stop 2>/dev/null || true'
	@sleep 1
	@echo "Current listeners on 5432 (if any):" 
	-@bash -c "ss -ltnp '( sport = :5432 )' 2>/dev/null || true"

docker-up: free-5432
	docker compose up -d db app

docker-down:
	docker compose down

wait-db:
	bash -c 'until pg_isready -h $${DB_HOST:-localhost} -p $${DB_PORT:-5432} -U $${DB_USER:-postgres}; do sleep 1; done'

migrate:
	psql "host=$${DB_HOST:-localhost} port=$${DB_PORT:-5432} user=$${DB_USER:-postgres} password=$${DB_PASSWORD:-postgres} dbname=$${DB_NAME:-r2_db} sslmode=$${DB_SSLMODE:-disable}" -f db/migrations/0001_init.sql

db-setup: docker-up wait-db migrate

## Apply all migrations (0001_*.sql, 0002_*.sql, ...)
migrate-all:
	@set -e; \
	for f in $$(ls -1 db/migrations/*.sql | sort); do \
		echo "Applying $$f"; \
		psql "host=$${DB_HOST:-localhost} port=$${DB_PORT:-5432} user=$${DB_USER:-postgres} password=$${DB_PASSWORD:-postgres} dbname=$${DB_NAME:-r2_db} sslmode=$${DB_SSLMODE:-disable}" -f $$f; \
	done

## Bring up DB+App and seed database
docker-up-seed: free-5432
	docker compose up -d db
	$(MAKE) wait-db
	$(MAKE) migrate-all
	docker compose up -d app
	@echo "Seed applied. Swagger: http://localhost:$${HTTP_PORT:-8080}/swagger"

# Build Docker image
docker-build:
	docker build -t r2-challenge:local .

# Run container (requires DB envs configured or external DB)
docker-run:
	docker run --rm -p 8080:8080 -e DB_HOST=host.docker.internal -e DB_USER=postgres -e DB_PASSWORD=postgres -e DB_NAME=r2_db r2-challenge:local

generate-mocks:
	./scripts/mock

swag-gen-app:
	# Generating swagger documentation
	swag init --parseDependency --parseInternal -ot yaml -g cmd/app/main.go --output cmd/app/swagger-gen
	# Creating temporary directory
	mkdir -p ./cmd/app/swagger-gen/tmp
	# Running generator CLI
	openapi-generator-cli generate -i ./cmd/app/swagger-gen/swagger.yaml -o ./cmd/app/swagger-gen/tmp -g openapi-yaml
	# Moving generated files from tmp to swagger-gen directory
	mv -f ./cmd/app/swagger-gen/tmp/openapi/openapi.yaml ./cmd/app/swagger-gen/swagger.yaml
	# Removing temporary directory
	rm -rf ./cmd/app/swagger-gen/tmp
	# Done

swaggen: swag-gen-app

# Rebuild and redeploy the app container with latest code and swagger
docker-rebuild:
	$(MAKE) swaggen
	docker compose build --no-cache app
	docker compose up -d --force-recreate app

# Rebuild and redeploy all services (app + db)
docker-rebuild-all:
	$(MAKE) swaggen
	docker compose build --no-cache
	docker compose up -d --force-recreate

# Rebuild all services and seed database (DB up, apply migrations, app up)
docker-rebuild-seed: free-5432
	$(MAKE) swaggen
	docker compose build --no-cache
	docker compose up -d db
	$(MAKE) wait-db
	$(MAKE) migrate-all
	docker compose up -d --force-recreate app
	@echo "Rebuilt app and applied seed. Swagger: http://localhost:$${HTTP_PORT:-8080}/swagger"

