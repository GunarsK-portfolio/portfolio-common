# Portfolio Common

![CI](https://github.com/GunarsK-portfolio/portfolio-common/workflows/CI/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/GunarsK-portfolio/portfolio-common)](https://goreportcard.com/report/github.com/GunarsK-portfolio/portfolio-common)
[![codecov](https://codecov.io/gh/GunarsK-portfolio/portfolio-common/branch/main/graph/badge.svg)](https://codecov.io/gh/GunarsK-portfolio/portfolio-common)

Shared Go package for common code across portfolio microservices.

## Overview

This module provides shared functionality used by multiple services in the portfolio application:

- Configuration management with validation
- Database models and repositories
- Authentication and security middleware
- Logging and metrics
- Common utilities and handlers

## Prerequisites

- Go 1.25+
- Node.js 22+ and npm 11+

## Packages

### `config`

Configuration management with validation and environment variable helpers.

#### ServiceConfig

Central service configuration with validation:

```go
type ServiceConfig struct {
    Port           string   // Server port (1-65535)
    Environment    string   // development, staging, or production
    AllowedOrigins []string // CORS allowed origins (required, no default)
}

cfg := config.NewServiceConfig("8080") // Default port if PORT not set
```

**Security Note**: `ALLOWED_ORIGINS` environment variable is required with no default. This forces explicit CORS configuration to prevent wildcard origin vulnerabilities.

#### DatabaseConfig

PostgreSQL connection configuration:

```go
type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
}

dbCfg := config.NewDatabaseConfig()
```

#### Other Configs

- `JWTConfig` - JWT secret and expiration settings
- `RedisConfig` - Redis connection configuration
- `S3Config` - MinIO/S3 storage configuration

#### Environment Helpers

- `GetEnv(key, defaultValue)` - Get env var with default
- `GetEnvRequired(key)` - Get required env var or panic
- `GetEnvBool(key, defaultValue)` - Get boolean env var
- `GetEnvInt(key, defaultValue)` - Get integer env var
- `GetEnvInt64(key, defaultValue)` - Get int64 env var

### `database`

Database connection utilities with connection pooling:

```go
import "github.com/GunarsK-portfolio/portfolio-common/database"

db, err := database.Connect(database.PostgresConfig{
    Host:     "localhost",
    Port:     "5432",
    User:     "portfolio_admin",
    Password: "password",
    DBName:   "portfolio",
})
```

### `middleware`

Gin middleware for authentication and security.

#### AuthMiddleware

JWT token validation via auth-service with automatic token TTL handling:

```go
import "github.com/GunarsK-portfolio/portfolio-common/middleware"

// Create auth middleware
authMiddleware := middleware.NewAuthMiddleware("http://auth-service:8084/api/v1")

// Apply to protected routes
protected := router.Group("/api/v1")
protected.Use(authMiddleware.ValidateToken())
protected.Use(authMiddleware.AddTTLHeader()) // Adds X-Token-TTL header
{
    protected.POST("/files", handler.Upload)
    protected.DELETE("/files/:id", handler.Delete)
}
```

**Optional timeout configuration:**

```go
authMiddleware := middleware.NewAuthMiddleware(
    "http://auth-service:8084/api/v1",
    middleware.WithTimeout(10 * time.Second),
)
```

#### SecurityMiddleware

CORS validation and security headers:

```go
import "github.com/GunarsK-portfolio/portfolio-common/middleware"

// Create security middleware with CORS configuration
securityMiddleware := middleware.NewSecurityMiddleware(
    cfg.AllowedOrigins,                    // []string from config
    "GET,POST,PUT,DELETE,OPTIONS",         // allowed methods
    "Content-Type,Authorization",          // allowed headers
    true,                                   // allow credentials
)

// Apply to all routes
router.Use(securityMiddleware.Apply())
```

**Features:**
- CORS origin whitelisting (no wildcard "*" support)
- Preflight request handling (OPTIONS)
- 403 response for disallowed origins
- Security headers: X-Content-Type-Options, X-Frame-Options, X-XSS-Protection
- Preflight caching (24 hours)
- Constructor validation (panics on empty/invalid origins)

**Service-specific examples:**

```go
// admin-api: Full CRUD operations
securityMiddleware := middleware.NewSecurityMiddleware(
    cfg.AllowedOrigins,
    "GET,POST,PUT,DELETE,OPTIONS",
    "Content-Type,Authorization",
    true,
)

// public-api: Read-only public access
securityMiddleware := middleware.NewSecurityMiddleware(
    cfg.AllowedOrigins,
    "GET,OPTIONS",
    "Content-Type",
    false,
)
```

### `models`

Shared database models:

- `Profile` - User profile information
- `WorkExperience` - Work history entries
- `Certification` - Professional certifications
- `Skill` - Technical skills with proficiency levels
- `Project` - Portfolio projects
- `Miniature` - Miniature painting showcase
- `StorageFile` - File metadata for MinIO storage

### `repository`

Shared database repository implementations.

#### ActionLogRepository

Audit logging for user actions:

```go
import "github.com/GunarsK-portfolio/portfolio-common/repository"

actionLogRepo := repository.NewActionLogRepository(db)

// Log a download
actionLogRepo.LogAction(&repository.ActionLog{
    ActionType:   "download",
    ResourceType: stringPtr("file"),
    ResourceID:   int64Ptr(fileID),
    UserID:       nil, // anonymous
    IPAddress:    stringPtr(clientIP),
    UserAgent:    stringPtr(userAgent),
})

// Get download count
count, err := actionLogRepo.CountActionsByResource("file", fileID)
```

### `handlers`

Common HTTP error handling utilities:

```go
import "github.com/GunarsK-portfolio/portfolio-common/handlers"

// Standardized error responses
handlers.RespondWithError(c, http.StatusBadRequest, "Invalid input")
handlers.RespondWithValidationError(c, validationErrors)
```

### `logger`

Structured logging with Gin middleware:

```go
import "github.com/GunarsK-portfolio/portfolio-common/logger"

// Initialize logger
logger.Init("production")

// Use logger middleware
router.Use(logger.GinLogger())

// Log messages
logger.Info("Server started", "port", 8080)
logger.Error("Database error", "error", err)
```

### `metrics`

Prometheus metrics collection:

```go
import "github.com/GunarsK-portfolio/portfolio-common/metrics"

metricsCollector := metrics.NewMetrics("admin-api")

// HTTP metrics are collected automatically via middleware
router.Use(metricsCollector.HTTPMetrics())

// Custom metrics
metricsCollector.IncrementCounter("files_uploaded", 1)
metricsCollector.ObserveHistogram("processing_time", duration)
```

### `utils`

Common utility functions:

```go
import "github.com/GunarsK-portfolio/portfolio-common/utils"

// Generate file URLs for MinIO storage
fileURL := utils.GenerateFileURL(baseURL, bucketName, objectKey)
```

## Development

### Available Commands

Using Task:

```bash
# Build and test
task build               # Build all packages (verify compilation)
task test                # Run tests
task test:coverage       # Run tests with coverage report
task clean               # Clean build artifacts

# Code quality
task format              # Format code with gofmt
task tidy                # Tidy and verify go.mod
task lint                # Run golangci-lint
task vet                 # Run go vet

# Security
task security:scan       # Run gosec security scanner
task security:vuln       # Check for vulnerabilities with govulncheck

# Development tools
task dev:install-tools   # Install dev tools (golangci-lint, govulncheck, etc.)

# CI/CD
task ci:all              # Run all CI checks
```

Using Go directly:

```bash
go build ./...       # Build all packages
go test ./...        # Run tests
go mod tidy          # Tidy dependencies
go mod verify        # Verify dependencies
```

## Services Using This Module

- `admin-api` - Admin dashboard backend (CRUD operations with auth)
- `auth-service` - Authentication and session management
- `files-api` - File upload/download with MinIO storage
- `public-api` - Public read-only portfolio data

## Version

This module follows semantic versioning. Current version: `v0.12.0`

## Breaking Changes

### v0.12.0

- **REQUIRED**: `ALLOWED_ORIGINS` environment variable must be set (comma-separated list of origins)
- No default value provided for security reasons
- Services will panic on startup if `ALLOWED_ORIGINS` is not configured

Example `.env`:
```
ALLOWED_ORIGINS=http://localhost:8080,http://localhost:8081,https://portfolio.example.com
```
