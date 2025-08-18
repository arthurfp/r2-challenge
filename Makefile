run:
	go run ./cmd/app

build:
	go build ./...

tidy:
	go mod tidy

test:
	go test ./...

lint:
	./bin/golangci-lint run ./...

docker-up:
	docker compose up -d db

docker-down:
	docker compose down

wait-db:
	bash -c 'until pg_isready -h $${DB_HOST:-localhost} -p $${DB_PORT:-5432} -U $${DB_USER:-postgres}; do sleep 1; done'

migrate:
	psql "host=$${DB_HOST:-localhost} port=$${DB_PORT:-5432} user=$${DB_USER:-postgres} password=$${DB_PASSWORD:-postgres} dbname=$${DB_NAME:-r2_db} sslmode=$${DB_SSLMODE:-disable}" -f db/migrations/0001_init.sql

db-setup: docker-up wait-db migrate


# Build Docker image
docker-build:
	docker build -t r2-challenge:local .

# Run container (requires DB envs configured or external DB)
docker-run:
	docker run --rm -p 8080:8080 -e DB_HOST=host.docker.internal -e DB_USER=postgres -e DB_PASSWORD=postgres -e DB_NAME=r2_db r2-challenge:local

