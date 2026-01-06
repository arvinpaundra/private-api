# Private API - Quiz Management System

A comprehensive RESTful API for managing educational quizzes, assessments, and submissions built with Go, following Domain-Driven Design (DDD) principles and Clean Architecture.

[![Go Version](https://img.shields.io/badge/Go-1.24.5-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 📋 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [API Documentation](#api-documentation)
- [Deployment](#deployment)

## Overview

Private API is an enterprise-grade quiz management system designed for educational institutions. It provides a robust backend for creating, managing, and taking quizzes with features like multi-choice questions, automated grading, real-time submissions, and comprehensive analytics.

## Features

- **Authentication & Authorization**

  - JWT-based authentication with access & refresh tokens
  - Secure password hashing with bcrypt

- **Quiz Management**

  - Create and manage subjects (Math, Physics, etc.)
  - Grade level organization (Grade 10, 11, 12, etc.)
  - Module creation with multiple-choice questions
  - Question validation (2-4 choices, single correct answer)
  - Publish/unpublish module toggle

- **Submission System**

  - Public quiz-taking interface
  - Real-time answer submission tracking
  - Automatic grading and scoring
  - Submission finalization with results

- **Dashboard & Analytics**

  - Real-time statistics
  - Module, subject, and grade counts
  - Submission tracking
  - User activity monitoring

- **Infrastructure**
  - Structured logging with Zap
  - Graceful shutdown
  - CORS support
  - Request validation
  - Error handling middleware

## Tech Stack

### Core

- **Language**: Go 1.24.5
- **Web Framework**: Gin (v1.10.1)
- **Database ORM**: GORM (v1.26.1)
- **CLI**: Cobra (v1.9.1)
- **Configuration**: Viper (v1.20.1)

### Database

- **Relational DB**: PostgreSQL 15+

### Security & Validation

- **JWT**: golang-jwt/jwt (v5.2.3)
- **Password Hashing**: bcrypt (golang.org/x/crypto)
- **Validation**: go-playground/validator (v10.26.0)

### Observability

- **Logging**: Zap (v1.27.0)
- **Health Checks**: Custom implementation

### DevOps & Deployment

- **Containerization**: Docker & Docker Compose
- **Orchestration**: Kubernetes
- **CI/CD**: GitHub Actions
- **Ingress**: Traefik
- **DNS**: External DNS with Cloudflare
- **Migrations**: golang-migrate

## Architecture

This application follows **Domain-Driven Design (DDD)** and **Clean Architecture** principles:

```
┌─────────────────────────────────────────────────────────┐
│                   Presentation Layer                     │
│              (REST Handlers, Middleware)                 │
├─────────────────────────────────────────────────────────┤
│                   Application Layer                      │
│              (Use Cases, Services, DTOs)                 │
├─────────────────────────────────────────────────────────┤
│                     Domain Layer                         │
│    (Entities, Value Objects, Domain Services)           │
├─────────────────────────────────────────────────────────┤
│                 Infrastructure Layer                     │
│   (Repositories, External Services, Database)           │
└─────────────────────────────────────────────────────────┘
```

### Domain Bounded Contexts

1. **Authentication Domain** - User management, login, sessions
2. **Subject Domain** - Educational subject management
3. **Grade Domain** - Grade level management
4. **Module Domain** - Quiz/exam content with questions
5. **Submission Domain** - Quiz-taking and answer submissions
6. **Dashboard Domain** - Analytics and statistics

### Key Patterns

- **Repository Pattern**: Data access abstraction
- **Unit of Work**: Transaction management
- **ACL (Anti-Corruption Layer)**: Domain isolation
- **Dependency Injection**: Loose coupling
- **Command/Query Separation**: Clear intent in operations

## Project Structure

```
.
├── application/           # Application layer (REST API)
│   └── rest/
│       ├── handler/       # HTTP request handlers
│       ├── middleware/    # Authentication, CORS, logging
│       └── router/        # Route definitions
├── cmd/                   # CLI commands
│   ├── rest.go           # REST server command
│   └── root.go           # Root command
├── config/               # Configuration management
│   └── config.go         # Viper + environment variables
├── core/                 # Shared utilities
│   ├── format/           # Response formatting
│   ├── token/            # JWT implementation
│   ├── trait/            # Common interfaces
│   ├── util/             # Utilities (logger, hash, uuid)
│   └── validator/        # Request validation
├── database/             # Database connections
│   ├── memorydb/         # Redis connection
│   └── relationaldb/     # PostgreSQL connection
├── domain/               # Domain layer (business logic)
│   ├── auth/             # Authentication domain
│   ├── dashboard/        # Dashboard domain
│   ├── grade/            # Grade domain
│   ├── module/           # Module domain
│   ├── subject/          # Subject domain
│   └── submission/       # Submission domain
│       ├── entity/       # Domain entities
│       ├── repository/   # Repository interfaces
│       ├── service/      # Domain services (use cases)
│       ├── constant/     # Domain constants & errors
│       └── response/     # Response DTOs
├── infrastructure/       # Infrastructure layer
│   ├── auth/            # Auth repositories & ACL
│   ├── dashboard/       # Dashboard ACL adapters
│   ├── grade/           # Grade repositories
│   ├── module/          # Module repositories & ACL
│   ├── subject/         # Subject repositories
│   └── submission/      # Submission repositories & ACL
├── model/               # Database models (GORM)
├── migrations/          # Database migrations
├── deploy/              # Deployment configurations
│   └── k8s/            # Kubernetes manifests
├── docs/
│   └── openapi/        # OpenAPI 3.0 specification
├── .github/
│   └── workflows/      # CI/CD pipelines
├── docker-compose.yaml
├── Dockerfile
├── Makefile
└── main.go
```

## Prerequisites

- **Go**: 1.24.5 or higher
- **Docker**: 20.10+ (for containerized development)
- **Docker Compose**: v2.0+ (for local stack)
- **PostgreSQL**: 15+ (if running locally)
- **Redis**: 7+ (if running locally)
- **golang-migrate**: For database migrations

## Installation

### 1. Clone the repository

```bash
git clone https://github.com/arvinpaundra/private-api.git
cd private-api
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Set up environment variables

Create a `.env` file in the root directory:

```bash
cp .env.example .env
```

Edit `.env` with your configuration (see [Configuration](#configuration) section).

## Configuration

Create a `.env` file with the following variables:

```env
# Application
APP_ENV=development          # development | production

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=root
DB_PASS=root
DB_DBNAME=private_db
DB_SSLMODE=disable          # disable | require | verify-full
DB_TIMEZONE=Asia/Jakarta

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASS=
REDIS_DB=0

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
```

**⚠️ Security Note**: Never commit `.env` files to version control. Always use strong, unique secrets in production.

## Running the Application

### Option 1: Using Docker Compose (Recommended)

```bash
# Start all services (API, PostgreSQL, Redis)
docker-compose up -d

# View logs
docker-compose logs -f rest

# Stop services
docker-compose down
```

The API will be available at `http://localhost:80`

### Option 2: Using Make (Local Development)

```bash
# Run migrations
make migrateup

# Start REST server (default port: 8000)
make rest

# Or specify custom port
make rest REST_PORT=9000
```

### Option 3: Direct Go Run

```bash
# Run migrations first
migrate -path ./migrations -database "postgres://root:root@localhost:5432/private_db?sslmode=disable" up

# Start server
go run main.go rest -p 8000
```

## API Documentation

### OpenAPI/Swagger Documentation

The API follows OpenAPI 3.0 specification. The complete API documentation is available in:

```
docs/openapi/openapi.yaml
```

### View API Documentation

#### Option 1: Using Swagger Editor (Online)

1. Visit [Swagger Editor](https://editor.swagger.io/)
2. Copy the content of `docs/openapi/openapi.yaml`
3. Paste into the editor

#### Option 2: Using Swagger UI (Docker)

```bash
# Run Swagger UI container
docker run -p 8080:8080 \
  -e SWAGGER_JSON=/openapi.yaml \
  -v $(pwd)/docs/openapi/openapi.yaml:/openapi.yaml \
  swaggerapi/swagger-ui

# Access at http://localhost:8080
```

#### Option 3: Using VS Code Extension

1. Install "OpenAPI (Swagger) Editor" extension
2. Open `docs/openapi/openapi.yaml`
3. Right-click → "OpenAPI: Show Preview"

### API Endpoints Overview

```
Authentication
  POST   /v1/auth/register          - Register new user
  POST   /v1/auth/login             - User login
  POST   /v1/auth/logout            - User logout
  POST   /v1/auth/refresh           - Refresh access token

Dashboard (Protected)
  GET    /v1/dashboard/statistics   - Get system statistics

Subjects (Protected)
  POST   /v1/subjects               - Create subject
  GET    /v1/subjects               - List subjects
  GET    /v1/subjects/:id           - Get subject details
  PUT    /v1/subjects/:id           - Update subject
  DELETE /v1/subjects/:id           - Delete subject

Grades (Protected)
  POST   /v1/grades                 - Create grade
  GET    /v1/grades                 - List grades
  GET    /v1/grades/:id             - Get grade details
  PUT    /v1/grades/:id             - Update grade
  DELETE /v1/grades/:id             - Delete grade

Modules (Protected)
  POST   /v1/modules                          - Create module
  GET    /v1/modules                          - List modules
  GET    /v1/modules/:slug                    - Get module details
  GET    /v1/modules/:slug/questions          - Get module questions
  POST   /v1/modules/:slug/questions          - Add questions
  PATCH  /v1/modules/:slug/publish            - Toggle publish status
  DELETE /v1/modules/:slug                    - Delete module

Modules (Public)
  GET    /v1/modules/:slug/published                     - Get published module details
  GET    /v1/modules/:slug/questions/:question_slug      - Get published question

Submissions (Public)
  POST   /v1/modules/:slug/submissions                     - Start submission
  POST   /v1/modules/:slug/submissions/:code/answers       - Submit answer
  PATCH  /v1/modules/:slug/submissions/:code/finalize      - Finalize submission

Submissions (Protected)
  GET    /v1/submissions            - List all submissions
```

### Authentication

Protected endpoints require a Bearer token:

```bash
curl -H "Authorization: Bearer <access_token>" \
  http://localhost:8000/v1/subjects
```

## Database Migrations

### Create new migration

```bash
make migrateadd NAME=create_table_users
```

This creates two files:

- `migrations/<timestamp>_create_table_users.up.sql`
- `migrations/<timestamp>_create_table_users.down.sql`

### Run migrations

```bash
# Apply all pending migrations
make migrateup

# Rollback last migration
make migratedown

# Rollback and reapply (useful for development)
make migraterefresh

# Migrate to specific version
make migrateto VERSION=20251220154919
```

### Manual migration

```bash
# Set database URL
export DB_URL="postgres://user:pass@localhost:5432/dbname?sslmode=disable"

# Run migrations
migrate -path ./migrations -database "$DB_URL" up
```

## Deployment

### Docker Build

```bash
# Build image
docker build -t arvinpaundra/private-api:latest .

# Run container
docker run -p 80:80 \
  -e DB_HOST=your-db-host \
  -e JWT_SECRET=your-secret \
  arvinpaundra/private-api:latest
```

### Kubernetes Deployment

The application includes complete Kubernetes manifests in `deploy/k8s/`:

```bash
# Apply namespace
kubectl apply -f deploy/k8s/namespace.yaml

# Apply configurations
kubectl apply -f deploy/k8s/configmap.yaml

# Create secrets (update values first)
kubectl apply -f deploy/k8s/secret.yaml

# Deploy application
kubectl apply -f deploy/k8s/deployment.yaml

# Create service
kubectl apply -f deploy/k8s/service.yaml

# Setup ingress (with Traefik)
kubectl apply -f deploy/k8s/ingress.yaml
```

### CI/CD with GitHub Actions

The repository includes a GitHub Actions workflow (`.github/workflows/deploy.yaml`) that:

1. Builds Docker image on tag push (`v*.*.*`)
2. Pushes to Docker Hub
3. Deploys to Kubernetes cluster
4. Performs rolling update with zero downtime

**To deploy:**

```bash
# Create and push a tag
git tag v1.0.0
git push origin v1.0.0

# GitHub Actions will automatically deploy
```

### Environment-specific Configuration

The application reads configuration from:

1. `.env` file (local development)
2. Environment variables (Docker/Kubernetes)
3. Kubernetes Secrets and ConfigMaps (production)

Priority: **Environment Variables** > **Viper Config** (automatic fallback)

## Acknowledgements

- Go community for excellent libraries
- Domain-Driven Design principles by Eric Evans
- Clean Architecture by Robert C. Martin

---

**Happy Coding! 🚀**

For questions or issues, please [open an issue](https://github.com/arvinpaundra/private-api/issues) on GitHub.
