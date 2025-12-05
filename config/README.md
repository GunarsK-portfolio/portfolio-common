# config

Configuration management with validation and environment variable helpers.

## Usage

```go
import "github.com/GunarsK-portfolio/portfolio-common/config"

// Service configuration
cfg := config.NewServiceConfig("8080")

// Database configuration
dbCfg := config.NewDatabaseConfig()

// Environment helpers
value := config.GetEnv("KEY", "default")
required := config.GetEnvRequired("KEY")
boolVal := config.GetEnvBool("FEATURE_FLAG", false)
intVal := config.GetEnvInt("PORT", 8080)
```

## Configuration Types

- `ServiceConfig` - Port, environment, allowed origins
- `DatabaseConfig` - PostgreSQL connection settings
- `JWTConfig` - JWT secret and expiration
- `RedisConfig` - Redis connection settings
- `S3Config` - MinIO/S3 storage settings
- `RabbitMQConfig` - RabbitMQ connection and queue settings
- `CookieConfig` - httpOnly cookie settings for authentication

## Environment Variables

### ServiceConfig

- `PORT` - Server port (default from constructor)
- `ENVIRONMENT` - development, staging, production
- `ALLOWED_ORIGINS` - CORS origins (required, comma-separated)

### DatabaseConfig

- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` - Required
- `DB_SSLMODE` - Optional (disable, require, verify-ca, verify-full)

### CookieConfig

- `COOKIE_DOMAIN` - Cookie domain (e.g., ".example.com" for prod, "" for local)
- `COOKIE_SECURE` - true for HTTPS only (default: false)
- `COOKIE_SAMESITE` - Strict, Lax, or None (default: Lax)
- `COOKIE_PATH` - Access token cookie path (default: "/")
- `COOKIE_REFRESH_PATH` - Refresh token cookie path (default: "/"). Must match
  the refresh endpoint URL as seen by the browser (e.g., "/auth/v1/refresh" for
  dev, "/api/v1/auth/refresh" for prod)
