# Changelog

## v0.34.0

**BREAKING**: JWT token generation methods now require a `scopes` parameter.

```go
// Before (v0.33.0)
token, err := jwtService.GenerateAccessToken(userID, username)
refreshToken, err := jwtService.GenerateRefreshToken(userID, username)

// After (v0.34.0)
scopes := map[string]string{"profile": "read", "projects": "edit"}
token, err := jwtService.GenerateAccessToken(userID, username, scopes)
refreshToken, err := jwtService.GenerateRefreshToken(userID, username, scopes)

// For nil scopes (no permissions)
token, err := jwtService.GenerateAccessToken(userID, username, nil)
```

- `GenerateAccessToken(userID, username)` now requires third `scopes` param
- `GenerateRefreshToken(userID, username)` now requires third `scopes` param
- Add `Scopes` field to JWT `Claims` struct
- Add `middleware/permission.go` with `RequirePermission()` middleware
- Add permission level constants: `LevelNone`, `LevelRead`, `LevelEdit`, `LevelDelete`
- Auth middleware now extracts scopes from JWT into Gin context

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
