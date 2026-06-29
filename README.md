# Inventory API

Simple REST API for inventory management built with Go, PostgreSQL, and GORM.

## Features

- Full CRUD operations for items
- PostgreSQL database with GORM ORM
- Configuration via `.env`
- Docker support
- GitHub Actions CI/CD
- Unit tests with 75%+ coverage

## Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.22 |
| Database | PostgreSQL 16 |
| ORM | GORM v2 |
| Config | godotenv |
| Container | Docker + Docker Compose |
| CI/CD | GitHub Actions |

## Project Structure

```
├── cmd/api/                # Entry point
├── internal/
│   ├── config/             # Configuration
│   ├── database/           # DB connection
│   ├── handler/            # HTTP handlers
│   ├── model/              # Data models
│   ├── repository/         # DB operations
│   └── service/            # Business logic
├── Dockerfile
├── docker-compose.yaml
└── .env.example
```

## Getting Started

### Prerequisites

- Go 1.22+
- Docker & Docker Compose (optional)
- PostgreSQL (if not using Docker)

### Quick Start with Docker

1. Clone repository

```bash
git clone https://github.com/afiffaizun/inventory-api.git
cd inventory-api
```

2. Copy environment file

```bash
cp .env.example .env
```

3. Run with Docker Compose

```bash
docker-compose up -d
```

4. Test the API

```bash
curl http://localhost:8080/health
```

### Manual Setup

1. Start PostgreSQL and create database `inventory_db`

2. Copy and configure `.env`

```bash
cp .env.example .env
```

3. Run the application

```bash
go run ./cmd/api
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | App info |
| GET | `/health` | Health check |
| GET | `/version` | Version info |
| GET | `/items` | List all items |
| GET | `/items/{id}` | Get item by ID |
| POST | `/items` | Create item |
| PUT | `/items/{id}` | Update item |
| DELETE | `/items/{id}` | Delete item |

### Examples

**Create Item**

```bash
curl -X POST http://localhost:8080/items \
  -H "Content-Type: application/json" \
  -d '{"code":"ITEM001","name":"Laptop","stock":10,"location":"Gudang A"}'
```

**Get All Items**

```bash
curl http://localhost:8080/items
```

**Get Item by ID**

```bash
curl http://localhost:8080/items/1
```

**Update Item**

```bash
curl -X PUT http://localhost:8080/items/1 \
  -H "Content-Type: application/json" \
  -d '{"code":"ITEM001","name":"Laptop Updated","stock":5}'
```

**Delete Item**

```bash
curl -X DELETE http://localhost:8080/items/1
```

## Testing

```bash
# Run all tests
go test ./... -v

# Run with coverage
go test ./... -cover

# Run specific package
go test ./internal/handler/ -v
```

## Development

```bash
# Build binary
go build -o inventory-api ./cmd/api

# Run linter
golangci-lint run

# Run locally without Docker
go run ./cmd/api
```

## CI/CD

GitHub Actions pipeline runs on push/PR:

1. **Lint** - golangci-lint
2. **Test** - go test with PostgreSQL
3. **Build** - go build
4. **Docker** - build & push to Docker Hub (main branch only)

## License

MIT
