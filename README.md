# go-users-api

A production-ready RESTful API built with **Go**, **GoFiber**, **PostgreSQL**, **SQLC**, **Uber Zap**, and **go-playground/validator**.

---

## Features

| Feature | Details |
|---|---|
| CRUD API | Create, Read, Update, Delete users |
| Dynamic Age | Age calculated at read time via Go's `time` package |
| SQLC | Type-safe, generated DB access layer |
| Validation | Input validation with `go-playground/validator` |
| Logging | Structured JSON logging with **Uber Zap** |
| Pagination | `GET /users?page=1&limit=10` |
| Middleware | `X-Request-ID` header + request duration logging |
| Docker | Multi-stage Dockerfile + Docker Compose |
| Tests | Unit tests for age calculation |

---

## Project Structure

```
.
├── cmd/server/main.go          # Entry point
├── config/config.go            # Env-based config
├── db/
│   ├── migrations/             # SQL migration files
│   └── sqlc/                   # SQLC-generated code
├── internal/
│   ├── handler/                # HTTP handlers
│   ├── repository/             # DB access layer
│   ├── service/                # Business logic
│   ├── routes/                 # Route registration
│   ├── middleware/             # RequestID + logger middleware
│   ├── models/                 # Request/response DTOs
│   └── logger/                 # Zap logger initializer
├── query.sql                   # SQLC query definitions
├── sqlc.yaml                   # SQLC configuration
├── Dockerfile                  # Multi-stage build
├── docker-compose.yml          # App + PostgreSQL
└── .env.example                # Environment variable template
```

---

## Prerequisites

- Go 1.22+
- PostgreSQL 14+ **or** Docker + Docker Compose

---

## Quick Start — Docker (Recommended)

```bash
# 1. Clone the repo
git clone <your-repo-url>
cd go-users-api

# 2. Copy env file (optional — defaults work out of the box)
cp .env.example .env

# 3. Start PostgreSQL + API
docker-compose up --build
```

The API will be available at **http://localhost:3000**.

---

## Quick Start — Local

```bash
# 1. Start PostgreSQL (skip if already running)
docker run -d \
  --name go_users_postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=usersdb \
  -p 5432:5432 \
  postgres:16-alpine

# 2. Apply the migration
psql -h localhost -U postgres -d usersdb -f db/migrations/000001_create_users.sql

# 3. Copy and edit env
cp .env.example .env

# 4. Run the server
go run ./cmd/server
```

---

## API Reference

### Health Check

```http
GET /health
```

### Create User

```http
POST /users
Content-Type: application/json

{
  "name": "Alice",
  "dob": "1990-05-10"
}
```

**Response 201**
```json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10"
}
```

---

### Get User by ID

```http
GET /users/:id
```

**Response 200**
```json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10",
  "age": 35
}
```

---

### Update User

```http
PUT /users/:id
Content-Type: application/json

{
  "name": "Alice Updated",
  "dob": "1991-03-15"
}
```

**Response 200**
```json
{
  "id": 1,
  "name": "Alice Updated",
  "dob": "1991-03-15"
}
```

---

### Delete User

```http
DELETE /users/:id
```

**Response 204 No Content**

---

### List All Users (Paginated)

```http
GET /users?page=1&limit=10
```

**Response 200**
```json
{
  "data": [
    {
      "id": 1,
      "name": "Alice",
      "dob": "1990-05-10",
      "age": 35
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total_items": 1,
    "total_pages": 1
  }
}
```

---

## HTTP Status Codes

| Code | Meaning |
|---|---|
| 200 | OK |
| 201 | Created |
| 204 | No Content (delete) |
| 400 | Bad Request (validation error) |
| 404 | User Not Found |
| 500 | Internal Server Error |

---

## Running Tests

```bash
go test ./internal/service/... -v
```

---

## Regenerating SQLC Code

If you modify `query.sql` or the migration schema:

```bash
# Install sqlc
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Regenerate
sqlc generate
```

---

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `SERVER_PORT` | `3000` | HTTP server port |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | PostgreSQL user |
| `DB_PASSWORD` | `postgres` | PostgreSQL password |
| `DB_NAME` | `usersdb` | PostgreSQL database name |
| `DB_SSLMODE` | `disable` | PostgreSQL SSL mode |
