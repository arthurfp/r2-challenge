run:
	go run ./cmd/app

build:
	go build ./...

tidy:
	go mod tidy

test:
	go test ./...

docker-up:
	docker compose up -d db

docker-down:
	docker compose down

wait-db:
	bash -c 'until pg_isready -h $${DB_HOST:-localhost} -p $${DB_PORT:-5432} -U $${DB_USER:-postgres}; do sleep 1; done'

migrate:
	psql "host=$${DB_HOST:-localhost} port=$${DB_PORT:-5432} user=$${DB_USER:-postgres} password=$${DB_PASSWORD:-postgres} dbname=$${DB_NAME:-r2_db} sslmode=$${DB_SSLMODE:-disable}" -f db/migrations/0001_init.sql

db-setup: docker-up wait-db migrate

