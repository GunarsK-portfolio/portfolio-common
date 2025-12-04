# middleware

Gin middleware for authentication and security.

## AuthMiddleware

Local JWT token validation with automatic TTL handling:

```go
import (
    "github.com/GunarsK-portfolio/portfolio-common/jwt"
    "github.com/GunarsK-portfolio/portfolio-common/middleware"
)

jwtService, _ := jwt.NewValidatorOnly(secret)
authMiddleware := middleware.NewAuthMiddleware(jwtService)

protected := router.Group("/api/v1")
protected.Use(authMiddleware.ValidateToken())
protected.Use(authMiddleware.AddTTLHeader())
```

### Token Extraction

The middleware extracts JWT tokens in the following order:

1. **Cookie** (httpOnly `access_token` cookie) - for browser requests
2. **Authorization header** (`Bearer <token>`) - for service-to-service calls

This allows secure httpOnly cookie authentication for browsers while maintaining
backwards compatibility with Authorization header for API clients.

After validation, access user info in handlers:

```go
userID, _ := c.Get("user_id")    // int64
username, _ := c.Get("username") // string
ttl, _ := c.Get("token_ttl")     // int64
```

## SecurityMiddleware

CORS validation and security headers:

```go
securityMiddleware := middleware.NewSecurityMiddleware(
    allowedOrigins,  // []string
    "GET,POST,PUT,DELETE,OPTIONS",
    "Content-Type,Authorization",
    true, // allow credentials
)
router.Use(securityMiddleware.Apply())
```

Features: CORS whitelisting, preflight handling, security headers
(X-Content-Type-Options, X-Frame-Options, X-XSS-Protection).
