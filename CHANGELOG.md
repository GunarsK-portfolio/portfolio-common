# Changelog

## v0.33.0

- Add `health` package for dependency health checking
- Add `Connection()` method to RabbitMQPublisher for health checks
- Add per-package README files
- Restructure main README to link to package docs

## v0.32.0

- Add `CloseDB` helper function to database package
- Add `queue` package with RabbitMQ publisher, retry queues, and DLQ support

## v0.21.0

**BREAKING**: AuthMiddleware now uses local JWT validation.

```go
// Before (v0.20.0)
authMiddleware := middleware.NewAuthMiddleware("http://auth-service:8084/api/v1")

// After (v0.21.0)
jwtService, _ := jwt.NewValidatorOnly(os.Getenv("JWT_SECRET"))
authMiddleware := middleware.NewAuthMiddleware(jwtService)
```

- `NewAuthMiddleware(authServiceURL string, opts...)` changed to
  `NewAuthMiddleware(jwtService jwt.Service)`
- Add `jwt` package for local token validation
- Remove `WithTimeout` option (no longer needed)
- Services must provide `JWT_SECRET` environment variable

## v0.19.0

- Add `SSLMode` field to `DatabaseConfig` with validation
- New environment variable: `DB_SSLMODE` (optional, default: `disable`)
- Valid values: `disable`, `allow`, `prefer`, `require`, `verify-ca`, `verify-full`

## v0.12.0

**BREAKING**: `ALLOWED_ORIGINS` environment variable is now required.

- No default value provided for security reasons
- Services will panic on startup if not configured
- Use comma-separated list: `ALLOWED_ORIGINS=http://localhost:8080,https://example.com`
