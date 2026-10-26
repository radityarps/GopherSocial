# GopherSocial 🦫

A production-grade RESTful social media API built with Go (Golang), PostgreSQL, and Redis. This project is a strict 1:1 implementation following the **[Complete Backend Engineering in Go](https://youtu.be/h3fqD6IprIA?si=UmHU0UTtVCp9lP9m)** course by Tiago Taquelim.

The project demonstrates production backend fundamentals: clean layered architecture, robust relational schema management, database connection pooling, stateless authentication, distributed caching, rate limiting, and automated API documentation.

---

## 🏗️ Architecture & Core Concepts

This application follows standard Go clean architecture practices, cleanly separating concerns between transport, business logic, and data persistence layers.

```mermaid
flowchart TD
    Client(["HTTP Client / Frontend"]) --> Router["Chi Router & Middleware<br/>(Auth, Rate Limiting, Logging)"]
    Router --> Handlers["Transport Layer<br/>(HTTP Handlers in cmd/api)"]
    Handlers --> Store["Repository / Store Layer<br/>(internal/store)"]
    Store --> Postgres[("PostgreSQL<br/>(Persistent Data)")]
    Handlers --> Cache[("Redis Cache<br/>(In-Memory Feed & User Cache)")]
    Handlers --> Mailer["Mailer Service<br/>(Account Activation & Notifications)"]
```

### Key Engineering Principles
- **Separation of Concerns**: HTTP handlers parse request payloads, validate parameters, and delegate data access entirely to the `store` layer.
- **Connection Pooling**: PostgreSQL connections are pooled (`MaxOpenConns`, `MaxIdleConns`, `ConnMaxLifetime`) to handle concurrent traffic without exhausting database resources.
- **Controlled Schema Evolution**: All database schema changes are strictly versioned through SQL migration scripts (`golang-migrate`), ensuring reproducibility across environments.
- **Stateless Authentication**: Security is managed via signed JSON Web Tokens (JWT) with salted bcrypt password hashing and Role-Based Access Control (RBAC).
- **Cache-Aside Pattern**: High-frequency queries (e.g. user feeds) leverage Redis for sub-millisecond retrieval with write-through invalidation.
- **Defensive API Design**: Token-bucket / Redis-backed rate limiting protects endpoints from abuse and brute-force traffic.

---

## 🛠️ Tech Stack

| Layer | Technology | Details |
| :--- | :--- | :--- |
| **Language** | [Go (Golang)](https://go.dev/) | Version 1.22+ |
| **HTTP Router** | [`go-chi/chi/v5`](https://github.com/go-chi/chi) | Idiomatic, standard-library compatible HTTP routing |
| **Primary Database** | [PostgreSQL 16](https://www.postgresql.org/) | Relational database with ACID transactions |
| **Database Migrations** | [`golang-migrate`](https://github.com/golang-migrate/migrate) | Versioned SQL up/down migrations |
| **Cache & Rate Limiting** | [Redis 7](https://redis.io/) | In-memory key-value store for caching and rate limiting |
| **Auth & Security** | `golang-jwt/jwt/v5` & `bcrypt` | Stateless bearer token authentication & password encryption |
| **Validation** | [`go-playground/validator/v10`](https://github.com/go-playground/validator) | Struct payload validation |
| **Documentation** | [Swagger / Swaggo](https://github.com/swaggo/swag) | Automated OpenAPI specification and Swagger UI |
| **Infrastructure** | [Docker Compose](https://docs.docker.com/compose/) | Containerized PostgreSQL and Redis services |

---

## 📁 Project Structure

```text
GopherSocial/
├── cmd/
│   └── api/                    # Application entrypoint & HTTP transport layer
│       ├── main.go             # Environment parsing, dependency wiring, server lifecycle
│       ├── api.go              # Application struct and router assembly
│       ├── auth.go             # Authentication handlers (login, register)
│       ├── users.go            # User management handlers
│       ├── posts.go            # Post CRUD handlers
│       ├── comments.go         # Comments handlers
│       ├── feed.go             # User feed handler
│       ├── health.go           # Healthcheck endpoint
│       └── middleware.go       # Auth, rate limiting, logging, and recovery middleware
├── internal/                   # Private library code (non-exportable)
│   ├── auth/                   # JWT generation, token validation, password hashing
│   ├── config/                 # Centralized environment configuration
│   ├── database/               # Database connection manager & pool configurations
│   ├── mailer/                 # Email delivery client & templates (SendGrid/Mailtrap)
│   ├── ratelimiter/            # Rate limiting implementation
│   └── store/                  # Data access layer (Repository pattern)
│       ├── storage.go          # Master Storage interface & factory
│       ├── users.go            # User SQL queries & data models
│       ├── posts.go            # Post SQL queries & data models
│       ├── comments.go         # Comment SQL queries & data models
│       └── followers.go        # Follower relationships & feed queries
├── migrations/                 # Versioned SQL migration files (.up.sql / .down.sql)
├── docs/                       # API documentation & architecture assets
├── docker-compose.yml          # Local container definitions (PostgreSQL, Redis)
├── Makefile                    # Automation shortcuts (migrate, build, run)
├── go.mod                      # Go module dependencies
└── README.md
```

---

## 🔌 API Endpoints (Core Overview)

| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :---: |
| `GET` | `/v1/health` | Service health status | No |
| `POST` | `/v1/authentication/user` | Register a new user | No |
| `POST` | `/v1/authentication/token` | Log in and receive JWT token | No |
| `PUT` | `/v1/users/activate/{token}` | Activate registered user account | No |
| `GET` | `/v1/users/{id}` | Retrieve user profile by ID | Yes |
| `PUT` | `/v1/users/{id}` | Update user profile | Yes (Owner) |
| `PUT` | `/v1/users/{id}/follow` | Follow a user | Yes |
| `PUT` | `/v1/users/{id}/unfollow` | Unfollow a user | Yes |
| `GET` | `/v1/users/feed` | Retrieve personalized social feed | Yes |
| `POST` | `/v1/posts` | Create a new post | Yes |
| `GET` | `/v1/posts/{id}` | Fetch a post by ID | Yes |
| `PATCH` | `/v1/posts/{id}` | Edit an existing post | Yes (Author) |
| `DELETE` | `/v1/posts/{id}` | Delete a post | Yes (Author/Admin) |
| `POST` | `/v1/posts/{id}/comments` | Add comment to a post | Yes |
| `GET` | `/swagger/*` | Swagger interactive documentation UI | No |

---

## 🚀 Getting Started

### Prerequisites
- **[Go](https://go.dev/dl/)** (1.22 or newer)
- **[Docker & Docker Compose](https://docs.docker.com/get-docker/)**
- **[golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)**

### 1. Environment Setup
Copy the sample environment file and configure variables:
```bash
cp .env.example .env
```

Key environment configuration variables:
```env
ADDR=:8080
DB_ADDR=postgres://postgres:postgres@localhost:5432/gophersocial?sslmode=disable
DB_MAX_OPEN_CONNS=30
DB_MAX_IDLE_CONNS=30
DB_MAX_IDLE_TIME=15m
REDIS_ADDR=localhost:6379
REDIS_PW=
REDIS_DB=0
AUTH_TOKEN_SECRET=your-secure-jwt-secret
AUTH_TOKEN_EXP=72h
MAILER_API_KEY=your-sendgrid-api-key
```

### 2. Launch Local Services
Start PostgreSQL and Redis containers via Docker Compose:
```bash
docker-compose up -d
```

### 3. Run Database Migrations
Apply versioned migrations to instantiate database tables and constraints:
```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/gophersocial?sslmode=disable" up
```

### 4. Generate Swagger Documentation (Optional)
Generate OpenAPI specifications from code annotations:
```bash
swag init -g ./cmd/api/main.go -o ./docs
```

### 5. Start the API Server
Run the HTTP server:
```bash
go run ./cmd/api/main.go
```
The server will boot on `http://localhost:8080`. Interactive API documentation can be accessed at `http://localhost:8080/swagger/index.html`.

---

## 🛠️ Helpful Commands

```bash
# Start Docker services
make docker-up

# Stop Docker services
make docker-down

# Run all migrations up
make migrate-up

# Rollback latest migration
make migrate-down

# Generate API docs
make gen-docs

# Run the API application
make run
```

---

## 📖 Reference & Credits

- Tutorial: **[Complete Backend Engineering in Go - Tiago Taquelim](https://youtu.be/h3fqD6IprIA?si=UmHU0UTtVCp9lP9m)**
- Original Project: **[sikozonpc/GopherSocial](https://github.com/sikozonpc/GopherSocial)**
