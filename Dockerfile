# syntax=docker/dockerfile:1

ARG GO_VERSION=1.23

FROM golang:${GO_VERSION}-alpine AS build
WORKDIR /src

# Enable Go modules and caching
ENV CGO_ENABLED=0

# Dependencies first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/app


FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=build /out/app /app/app
# Copy generated OpenAPI spec used by the app at runtime
COPY --from=build /src/cmd/app/swagger-gen/swagger.yaml /app/cmd/app/swagger-gen/swagger.yaml

EXPOSE 8080 9090

ENV HTTP_PORT=8080

USER nonroot
ENTRYPOINT ["/app/app"]

