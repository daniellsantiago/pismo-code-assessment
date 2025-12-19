# Pismo Code Assessment

REST API for managing customer accounts and transactions.

## Requirements

- [Docker](https://www.docker.com/) and Docker Compose
- [Go 1.24+](https://golang.org/) (for local development)

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

### Run Locally

```bash
# Start PostgreSQL
docker-compose up postgres -d

# Run the API
go run ./cmd/api
```

## API Endpoints

### Accounts

#### Create Account
```bash
curl -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -d '{"document_number": "12345678900"}'
```

Response:
```json
{
  "account_id": 1,
  "document_number": "12345678900"
}
```

#### Get Account
```bash
curl http://localhost:8080/accounts/1
```

Response:
```json
{
  "account_id": 1,
  "document_number": "12345678900"
}
```

### Transactions

#### Create Transaction
```bash
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{"account_id": 1, "operation_type_id": 4, "amount": 123.45}'
```

Response:
```json
{
  "transaction_id": 1,
  "account_id": 1,
  "operation_type_id": 4,
  "amount": 123.45
}
```

**Operation Types:**

| ID | Description | Amount |
|----|-------------|--------|
| 1 | PURCHASE | Negative |
| 2 | INSTALLMENT PURCHASE | Negative |
| 3 | WITHDRAWAL | Negative |
| 4 | PAYMENT | Positive |

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
├── cmd/
│   └── api/
│       └── main.go              # Application entrypoint
├── docs/
│   └── openapi.yaml             # API documentation
├── internal/
│   ├── domain/                  # Business entities and interfaces
│   ├── infrastructure/
│   │   ├── config/              # Configuration
│   │   ├── database/            # Repository implementations
│   │   │   └── migrations/      # SQL migrations
│   │   └── http/
│   │       ├── dto/             # Request/Response DTOs
│   │       ├── handler/         # HTTP handlers
│   │       ├── response/        # Shared response utilities
│   │       ├── router/          # Routes
│   │       └── server/          # HTTP server
│   └── usecase/                 # Application use cases
├── test/
│   └── integration/             # Integration tests
├── docker-compose.yaml
├── Dockerfile
└── README.md
```

## Tech Stack

- **Language:** Go 1.24
- **Database:** PostgreSQL 16
- **Container:** Docker
- **Documentation:** OpenAPI 3.0

