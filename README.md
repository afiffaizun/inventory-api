# Inventory System

Full-stack inventory management system with Go backend and Vue 3 frontend.

## Features

- Full CRUD operations for items
- PostgreSQL database with GORM ORM
- Vue 3 frontend with PrimeVue components
- Docker support
- GitHub Actions CI/CD
- Unit tests with 75%+ coverage

## Tech Stack

| Component | Technology |
|-----------|------------|
| **Backend** | |
| Language | Go 1.22 |
| Database | PostgreSQL 16 |
| ORM | GORM v2 |
| Config | godotenv |
| **Frontend** | |
| Framework | Vue 3 (Composition API) |
| Build | Vite |
| Language | TypeScript |
| CSS | Tailwind CSS |
| Components | PrimeVue 4 |
| **DevOps** | |
| Container | Docker + Docker Compose |
| CI/CD | GitHub Actions |

## Project Structure

```
├── backend/                 # Go API server
│   ├── cmd/api/             # Entry point
│   ├── internal/
│   │   ├── config/          # Configuration
│   │   ├── database/        # DB connection
│   │   ├── handler/         # HTTP handlers
│   │   ├── model/           # Data models
│   │   ├── repository/      # DB operations
│   │   └── service/         # Business logic
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile
├── frontend/                # Vue 3 app
│   ├── src/
│   │   ├── api/             # API service
│   │   ├── components/      # Vue components
│   │   ├── views/           # Page views
│   │   └── types/           # TypeScript types
│   ├── package.json
│   └── Dockerfile
├── docker-compose.yaml
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.22+
- Node.js 20+
- Docker & Docker Compose (optional)
- PostgreSQL (if not using Docker)

### Quick Start with Docker

1. Clone repository

```bash
git clone https://github.com/afiffaizun/inventory-api.git
cd inventory-api
```

2. Start all services

```bash
docker-compose up -d
```

3. Access the application

- Frontend: http://localhost:3000
- API: http://localhost:8080
- Adminer: http://localhost:8081

**Adminer Login Settings:**

| Field | Value |
|-------|-------|
| System | PostgreSQL |
| Server | `db` |
| Username | `postgres` |
| Password | `postgres` |
| Database | `inventory_db` |

### Development Setup

**Backend:**

```bash
cd backend
cp .env.example .env
go run ./cmd/api
```

**Frontend:**

```bash
cd frontend
npm install
npm run dev
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

## Development

**Backend:**

```bash
cd backend
go run ./cmd/api          # Run server
go test ./... -v          # Run tests
go build -o api ./cmd/api # Build binary
```

**Frontend:**

```bash
cd frontend
npm run dev               # Dev server
npm run build             # Production build
npm run lint              # Lint code
```

## CI/CD

GitHub Actions pipeline runs on push/PR:

1. **Lint** - golangci-lint
2. **Test** - go test with PostgreSQL
3. **Build** - go build
4. **Docker** - build & push to Docker Hub (main branch only)

## License

MIT
