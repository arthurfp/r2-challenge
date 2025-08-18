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

docker-up:
	docker compose up -d db app

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

