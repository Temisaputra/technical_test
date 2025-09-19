# War-Onk Service

A clean-architecture Go project for managing supplier and other modules with GORM, PostgreSQL, and Swagger API documentation.

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Running the Application](#running-the-application)
- [Database Migration](#database-migration)
- [Swagger Documentation](#swagger-documentation)
- [Generating Mocks](#generating-mocks)
- [Project Structure](#project-structure)

---

## Prerequisites

- Go 1.21+
- PostgreSQL 14+
- `mockgen` for generating Go mocks
- `make` for running predefined tasks

---

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/war-onk.git
cd war-onk
```

2. Install dependencies:

```bash
go mod tidy
```

3. Configure environment variables (e.g., `.env`) for database connection.

---

## Running the Application

To start the REST API server:

```bash
go run main.go rest
```

The server will start on `http://localhost:8085`.

---

## Database Migration

Run migrations using the CLI commands:

- **Migrate Up**:

```bash
go run main.go migrate up
```

- **Migrate Down**:

```bash
go run main.go migrate down
```

- **Fresh Migration**:

```bash
go run main.go migrate fresh
```

---

## Swagger Documentation

After starting the server, access the Swagger UI at:

```
http://localhost:8085/swagger/index.html
```

Use this to explore all endpoints including:

- Supplier Create
- Supplier Get All
- Supplier Get By ID
- Supplier Update
- Supplier Delete

---

## Generating Mocks

For unit testing, generate mocks using:

```bash
mockgen -source=internal/repository/supplier_repository.go -destination=./shared/mock/repository/repository_mock.go -package repository
```

You can also add a Makefile target for convenience:

```makefile
generate-mocks:
	mockgen -source=internal/repository/supplier_repository.go -destination=./shared/mock/repository/repository_mock.go -package repository
```

Run with:

```bash
make generate-mocks
```

---

## Project Structure

```
├── cmd/              # CLI commands
├── delivery/         # HTTP handlers, presenters, requests/responses
├── domain/           # Business entities and interfaces
├── db/               # Database implementation and repositories
├── internal/         # Internal packages
├── shared/mock/      # Generated mocks for testing
├── main.go           # Entry point
```

- **Repositories**: Handle DB transactions (optional, panic-safe).
- **Usecases**: Business logic. Use `WithTransaction(ctx, func(txCtx) error)` for atomic multi-repo operations.
- **Delivery**: REST API handlers with Swagger annotations.

---

## Testing

Run unit tests:

```bash
go test ./... -v
```

- Use GoMock for repository interfaces
- Use Convey for BDD-style testing

---

## Notes

- Transactions are optional and panic-safe.
- Supports multi-repository transactional operations.
- Error handling uses custom `TxError` type to distinguish commit, rollback, and operation errors.
