# Portfolio Common

![CI](https://github.com/GunarsK-portfolio/portfolio-common/workflows/CI/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/GunarsK-portfolio/portfolio-common)](https://goreportcard.com/report/github.com/GunarsK-portfolio/portfolio-common)
[![codecov](https://codecov.io/gh/GunarsK-portfolio/portfolio-common/graph/badge.svg)](https://codecov.io/gh/GunarsK-portfolio/portfolio-common)
[![CodeRabbit](https://img.shields.io/coderabbit/prs/github/GunarsK-portfolio/portfolio-common?label=CodeRabbit&color=2ea44f)](https://coderabbit.ai)

Shared Go package for common code across portfolio microservices.

## Prerequisites

- Go 1.25+
- Node.js 22+ and npm 11+

## Packages

| Package | Description |
|---------|-------------|
| [config](config/) | Configuration management and environment helpers |
| [database](database/) | PostgreSQL connection with GORM |
| [jwt](jwt/) | Local JWT validation and generation |
| [middleware](middleware/) | Auth and security middleware for Gin |
| [models](models/) | Shared GORM database models |
| [audit](audit/) | Security event logging |
| [repository](repository/) | Shared repository implementations |
| [handlers](handlers/) | Common HTTP handler utilities |
| [logger](logger/) | Structured logging with slog |
| [metrics](metrics/) | Prometheus metrics collection |
| [server](server/) | HTTP server with graceful shutdown |
| [utils](utils/) | Common utility functions |
| [queue](queue/) | RabbitMQ publishing with retry/DLQ |
| [health](health/) | Dependency health checking |

## Quick Start

```go
import (
    "github.com/GunarsK-portfolio/portfolio-common/config"
    "github.com/GunarsK-portfolio/portfolio-common/database"
    "github.com/GunarsK-portfolio/portfolio-common/health"
    "github.com/GunarsK-portfolio/portfolio-common/jwt"
    "github.com/GunarsK-portfolio/portfolio-common/middleware"
)

// Configuration
cfg := config.NewServiceConfig("8080")
dbCfg := config.NewDatabaseConfig()

// Database
db, _ := database.Connect(dbCfg)
defer database.CloseDB(db)

// Auth middleware
jwtService, _ := jwt.NewValidatorOnly(cfg.JWTSecret)
authMiddleware := middleware.NewAuthMiddleware(jwtService)

// Health checks
healthAgg := health.NewAggregator(3 * time.Second)
healthAgg.Register(health.NewPostgresChecker(db))
router.GET("/health", healthAgg.Handler())
```

## Development

```bash
task ci:all              # Run all CI checks
task test                # Run tests
task lint                # Run linter
task format              # Format code
```

## Services Using This Module

- `admin-api` - Admin dashboard backend
- `auth-service` - Authentication and sessions
- `files-api` - File upload/download with MinIO
- `messaging-api` - Contact form and message queue
- `public-api` - Public read-only portfolio data

## Version

Current version: `v0.33.0`

See [CHANGELOG.md](CHANGELOG.md) for breaking changes and migration guides.

## License

[MIT](LICENSE)
