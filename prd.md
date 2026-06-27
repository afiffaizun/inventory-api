# PRD - Inventory API Development Plan

## Current State

| Layer | Status |
|-------|--------|
| Model | `Item`, `Response` |
| Repository | In-memory slice (data hilang saat restart) |
| Service | `GetHome`, `GetVersion`, `GetItems`, `CreateItem` |
| Handler | HTTP basic, no middleware |
| Config | Hardcoded (port 8080) |
| Testing | None |

---

## Development Roadmap

### Phase 1: Core Infrastructure (Prioritas Tinggi)

#### 1.1 Database Integration
- PostgreSQL/MySQL dengan driver atau ORM (GORM/sqlx)
- Migration file (`migrations/`)
- Implementasi `ItemRepository` sesungguhnya (bukan in-memory)

#### 1.2 Configuration Management
- Tambah `.env` + `envconfig` atau `viper` untuk konfigurasi (port, DB URL, dll)
- Buat `config/config.go`

#### 1.3 Error Handling & Response
- Buat standard error response struct
- Konsisten error codes (400, 404, 500)

---

### Phase 2: API Enhancement

#### 2.1 HTTP Router
- Ganti `net/http` ke **chi** atau **gin** untuk routing lebih structured
- Route grouping: `/api/v1/items`

#### 2.2 Expanded CRUD
- `GET /items/:id` — Get single item
- `PUT /items/:id` — Update item
- `DELETE /items/:id` — Delete item
- Search/filter (`?location=Warehouse A`)

#### 2.3 Validation
- Validasi input pada handler/service (e.g., `Code` wajib diisi, `Stock` >= 0)
- Gunakan `go-playground/validator` atau custom validation

#### 2.4 Pagination
- Untuk `GET /items` — tambah `?page=1&limit=10`

---

### Phase 3: Security & Middleware

#### 3.1 Middleware
- Logging middleware
- CORS middleware
- Authentication/Authorization (JWT)

---

### Phase 4: Quality & Reliability

#### 4.1 Testing
- Unit test untuk repository dan service
- Integration test untuk handler (gunakan `httptest`)
- Coverage target: 70%+

#### 4.2 Logging & Observability
- Structured logging (`zerolog` / `zap`)
- Health check endpoint sudah ada, tambah `/ready` endpoint
- Metrics (Prometheus)

---

### Phase 5: Deployment & DevOps

#### 5.1 Docker
- `Dockerfile` untuk build image
- `docker-compose.yaml` untuk jalankan app + database

#### 5.2 CI/CD Pipeline
- GitHub Actions untuk lint, test, build, push image
- `.golangci-lint.yaml` untuk code quality

---

## API Endpoints (Target)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Home / App Info |
| GET | `/health` | Health Check |
| GET | `/version` | Version Info |
| GET | `/api/v1/items` | List Items (with pagination) |
| GET | `/api/v1/items/:id` | Get Single Item |
| POST | `/api/v1/items` | Create Item |
| PUT | `/api/v1/items/:id` | Update Item |
| DELETE | `/api/v1/items/:id` | Delete Item |

---

## Tech Stack (Target)

- **Language:** Go 1.22+
- **Router:** chi / gin
- **Database:** PostgreSQL
- **ORM/Driver:** GORM / sqlx
- **Config:** viper / envconfig
- **Logger:** zerolog / zap
- **Validator:** go-playground/validator
- **Testing:** testing + httptest
- **Container:** Docker + Docker Compose
- **CI/CD:** GitHub Actions
