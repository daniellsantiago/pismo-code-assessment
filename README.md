# Pismo Code Assessment

REST API for managing customer accounts and transactions.

## Requirements

- [Go 1.24+](https://golang.org/)
- [PostgreSQL 16+](https://www.postgresql.org/)
- [Docker](https://www.docker.com/) and Docker Compose (highly recommended)

## Getting Started

### Run with Docker (Recommended)

```bash
# Start all services (API, PostgreSQL, Swagger UI)
docker-compose up --build

# Or run in background
docker-compose up --build -d
```

**Available Services:**

| Service | URL |
|---------|-----|
| API | http://localhost:8080 |
| Swagger UI | http://localhost:8081 |
| PostgreSQL | localhost:5432 |

### Or run locally

```bash
# Start PostgreSQL
docker-compose up postgres -d

# Run the API
go run ./cmd/api
```

> You can also use your own PostgreSQL instance instead of Docker. Just set the `DATABASE_URL` environment variable with your connection string.

## API Documentation

Full API documentation is available via Swagger UI at http://localhost:8081 when running with Docker.

You can also view the OpenAPI spec directly at [`docs/openapi.yaml`](docs/openapi.yaml).

## Running Tests

```bash
# Run all unit tests
go test ./internal/... -v

# Run integration tests (requires Docker)
go test ./test/integration/... -v

# Run tests with coverage
go test ./internal/... -cover
```

## Project Structure

```
├── cmd/api/                     # Application entrypoint
├── docs/                        # OpenAPI documentation
├── internal/
│   ├── domain/                  # Business entities and interfaces
│   ├── infrastructure/
│   │   ├── config/              # Configuration
│   │   ├── database/            # Repository implementations and migrations
│   │   └── http/                # HTTP handlers, middleware, router, server
│   └── usecase/                 # Application use cases
├── pkg/logger/                  # Shared logger package
└── test/integration/            # Integration tests
```
