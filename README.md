# Portfolio Common

Shared Go package for common code across portfolio microservices.

## Overview

This module provides shared functionality used by multiple services in the portfolio application:
- Database models and repositories
- Common utilities
- Shared types

## Packages

### `config`

Environment variable helpers:
- `GetEnv()` - Get env var with default value
- `GetEnvRequired()` - Get required env var or panic
- `GetEnvBool()` - Get boolean env var
- `GetEnvInt()` - Get integer env var
- `GetEnvInt64()` - Get int64 env var

### `database`

Database connection utilities:
- `Connect()` - PostgreSQL connection with connection pooling

### `middleware`

Gin middleware:
- `AuthMiddleware` - JWT token validation via auth-service

### `repository`

Shared database repository implementations:
- `ActionLogRepository` - Audit logging for user actions (logins, downloads, uploads, etc.)

## Usage

Import in your service:

```go
import "github.com/GunarsK-portfolio/portfolio-common/repository"
```

### Database Connection Example

```go
import "github.com/GunarsK-portfolio/portfolio-common/database"

// Connect to PostgreSQL with connection pooling
db, err := database.Connect(database.PostgresConfig{
    Host:     "localhost",
    Port:     "5432",
    User:     "portfolio_admin",
    Password: "password",
    DBName:   "portfolio",
})
```

### Auth Middleware Example

```go
import "github.com/GunarsK-portfolio/portfolio-common/middleware"

// Setup Gin router with auth middleware
authMiddleware := middleware.NewAuthMiddleware("http://auth-service:8084/api/v1")
protected := router.Group("/api/v1")
protected.Use(authMiddleware.ValidateToken())
{
    protected.POST("/files", handler.Upload)
}
```

### Action Logging Example

```go
import "github.com/GunarsK-portfolio/portfolio-common/repository"

// Initialize repository
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

## Development

### Installing Dependencies

```bash
go mod download
```

### Running Tests

```bash
go test ./...
```

## Services Using This Module

- `auth-service` - Logs login/logout actions
- `files-api` - Logs file downloads, uploads, deletes
- `admin-api` - (future) Logs admin actions
- `public-api` - (future) Logs public file downloads

## Version

This module follows semantic versioning. Current version: `v0.1.0`
