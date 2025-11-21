# Portfolio Common

![CI](https://github.com/GunarsK-portfolio/portfolio-common/workflows/CI/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/GunarsK-portfolio/portfolio-common)](https://goreportcard.com/report/github.com/GunarsK-portfolio/portfolio-common)
[![codecov](https://codecov.io/gh/GunarsK-portfolio/portfolio-common/branch/main/graph/badge.svg)](https://codecov.io/gh/GunarsK-portfolio/portfolio-common)

Shared Go package for common code across portfolio microservices.

## Overview

This module provides shared functionality used by multiple services in the
portfolio application:

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

**Security Note**: `ALLOWED_ORIGINS` environment variable is required with
no default. This forces explicit CORS configuration to prevent wildcard
origin vulnerabilities.

#### DatabaseConfig

PostgreSQL connection configuration:

```go
type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    Name     string
    SSLMode  string // disable (local), require (AWS RDS), verify-ca, verify-full
}

dbCfg := config.NewDatabaseConfig()
```

**Environment Variables:**

- `DB_HOST` - PostgreSQL host (required)
- `DB_PORT` - PostgreSQL port (required)
- `DB_USER` - Database user (required)
- `DB_PASSWORD` - Database password (required)
- `DB_NAME` - Database name (required)
- `DB_SSLMODE` - SSL mode (optional, default: `disable`)
  - Valid values: `disable`, `allow`, `prefer`, `require`, `verify-ca`, `verify-full`
  - Local Docker: Use `disable`
  - AWS RDS/Aurora: Use `require`

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
    SSLMode:  "disable", // disable (local), require (AWS RDS)
    TimeZone: "UTC",
})
```

### `jwt`

Local JWT token validation and generation.

```go
import "github.com/GunarsK-portfolio/portfolio-common/jwt"

// For services that only validate tokens (admin-api, files-api)
jwtService, err := jwt.NewValidatorOnly(cfg.JWTSecret)

// For services that generate and validate tokens (auth-service)
jwtService, err := jwt.NewService(cfg.JWTSecret, 15*time.Minute, 168*time.Hour)

// Validate a token
claims, err := jwtService.ValidateToken(tokenString)
if err != nil {
    // Token invalid or expired
}

// Access claims
userID := claims.UserID
username := claims.Username
ttl := claims.GetTTL() // Remaining seconds until expiry
```

**Benefits of local validation:**

- No network calls (faster, more resilient)
- No single point of failure
- Industry standard approach (used by Netflix, Google, Stripe)

### `middleware`

Gin middleware for authentication and security.

#### AuthMiddleware

Local JWT token validation with automatic TTL handling and user context:

```go
import (
    "github.com/GunarsK-portfolio/portfolio-common/jwt"
    "github.com/GunarsK-portfolio/portfolio-common/middleware"
)

// Create JWT service for validation
jwtService, _ := jwt.NewValidatorOnly(cfg.JWTSecret)

// Create auth middleware with JWT service
authMiddleware := middleware.NewAuthMiddleware(jwtService)

// Apply to protected routes
protected := router.Group("/api/v1")
protected.Use(authMiddleware.ValidateToken())
protected.Use(authMiddleware.AddTTLHeader()) // Adds X-Token-TTL header
{
    protected.POST("/files", handler.Upload)
    protected.DELETE("/files/:id", handler.Delete)
}
```

**After validation, user claims are stored in Gin context:**

```go
// In your handlers, access authenticated user information
userID, _ := c.Get("user_id")     // int64
username, _ := c.Get("username")  // string
ttl, _ := c.Get("token_ttl")      // int64 (seconds until expiry)
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
- Security headers: X-Content-Type-Options, X-Frame-Options,
  X-XSS-Protection
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

### `audit`

Centralized security event logging with automatic context extraction.

#### Quick Start

```go
import "github.com/GunarsK-portfolio/portfolio-common/audit"

// 1. Add audit context middleware (extracts IP and user agent)
router.Use(audit.ContextMiddleware())

// 2. Add auth middleware (stores user_id in context)
router.Use(authMiddleware.ValidateToken())

// 3. In handlers, log security events with one line
audit.LogFromContext(c, actionLogRepo, audit.ActionLoginSuccess, nil, nil, map[string]interface{}{
    "username": username,
})
```

**Automatic context extraction:**

- Client IP (X-Forwarded-For → X-Real-IP → RemoteAddr)
- User-Agent header
- user_id (from auth middleware)

#### Action Types

Predefined constants for consistency:

```go
audit.ActionLoginSuccess         // Successful login
audit.ActionLoginFailure         // Failed login attempt
audit.ActionLogout               // User logout
audit.ActionTokenRefresh         // Token refresh
audit.ActionTokenValidation      // Token validation failure
audit.ActionFileUpload           // File uploaded
audit.ActionFileDownload         // File downloaded
audit.ActionFileDelete           // File deleted
```

#### Resource Types

```go
audit.ResourceTypeFile           // File resource
audit.ResourceTypeUser           // User resource
```

#### Usage Examples

**Login success (authenticated):**

```go
// user_id automatically extracted from context by LogFromContext
err := audit.LogFromContext(c, actionLogRepo,
    audit.ActionLoginSuccess, nil, nil, map[string]interface{}{
        "username": username,
    })
```

**Login failure (no user_id):**

```go
err := audit.LogFromContext(c, actionLogRepo,
    audit.ActionLoginFailure, nil, nil, map[string]interface{}{
        "username": attemptedUsername,
        "reason": "invalid_credentials",
    })
```

**File download (with resource):**

```go
resourceType := audit.ResourceTypeFile
err := audit.LogFromContext(c, actionLogRepo,
    audit.ActionFileDownload, &resourceType, &fileID,
    map[string]interface{}{
        "filename": fileName,
        "size": fileSize,
    })
```

**Explicit user_id (when not using auth middleware):**

```go
// Use LogAction when user_id is not in context
err := audit.LogAction(c, actionLogRepo,
    audit.ActionLoginSuccess, nil, nil, &userID,
    map[string]interface{}{
        "method": "api_key",
    })
```

#### Helper Functions

```go
// Get context values (set by ContextMiddleware and auth middleware)
clientIP := audit.GetClientIP(c)      // *string
userAgent := audit.GetUserAgent(c)    // *string
userID := audit.GetUserID(c)          // *int64
```

### `repository`

Shared database repository implementations.

#### ActionLogRepository

Low-level audit log repository (use `audit` package helpers instead):

```go
import "github.com/GunarsK-portfolio/portfolio-common/repository"

actionLogRepo := repository.NewActionLogRepository(db)

// Query logs
logs, err := actionLogRepo.GetActionsByType("login_success", 100)
logs, err := actionLogRepo.GetActionsByUser(userID, 50)
logs, err := actionLogRepo.GetActionsByResource("file", fileID)
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

This module follows semantic versioning. Current version: `v0.21.0`

## Breaking Changes

### v0.21.0

- **BREAKING**: `AuthMiddleware` now uses local JWT validation instead of HTTP
  calls
- `NewAuthMiddleware(authServiceURL string, opts...)` changed to
  `NewAuthMiddleware(jwtService jwt.Service)`
- **ADDED**: New `jwt` package for local token validation
- **REMOVED**: `WithTimeout` option (no longer needed - no HTTP calls)
- Services must now provide `JWT_SECRET` environment variable
- Eliminates network dependency on auth-service for token validation

Migration:

```go
// Before (v0.20.0)
authMiddleware := middleware.NewAuthMiddleware("http://auth-service:8084/api/v1")

// After (v0.21.0)
jwtService, _ := jwt.NewValidatorOnly(os.Getenv("JWT_SECRET"))
authMiddleware := middleware.NewAuthMiddleware(jwtService)
```

### v0.19.0

- **ADDED**: `SSLMode` field to `DatabaseConfig` with validation
- New environment variable: `DB_SSLMODE` (optional, default: `disable`)
- Valid values: `disable`, `allow`, `prefer`, `require`, `verify-ca`, `verify-full`
- Validation prevents invalid values at startup

### v0.12.0

- **REQUIRED**: `ALLOWED_ORIGINS` environment variable must be set
  (comma-separated list of origins)
- No default value provided for security reasons
- Services will panic on startup if `ALLOWED_ORIGINS` is not configured

Example `.env`:

```bash
ALLOWED_ORIGINS=http://localhost:8080,https://portfolio.example.com
DB_SSLMODE=disable  # Local Docker: disable, AWS RDS: require
```
