run:
	go run ./cmd/app

build:
	go build ./...

tidy:
	go mod tidy

test:
	go test ./...

