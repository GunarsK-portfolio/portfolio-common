# jwt

Local JWT token validation and generation using HS256 signing.

## Usage

```go
import "github.com/GunarsK-portfolio/portfolio-common/jwt"

// For services that only validate tokens
jwtService, err := jwt.NewValidatorOnly(secret)

// For services that generate and validate tokens
jwtService, err := jwt.NewService(secret, 15*time.Minute, 168*time.Hour)

// Validate a token
claims, err := jwtService.ValidateToken(tokenString)
userID := claims.UserID
username := claims.Username
ttl := claims.GetTTL()

// Generate tokens (full service only)
accessToken, err := jwtService.GenerateAccessToken(userID, username)
refreshToken, err := jwtService.GenerateRefreshToken(userID, username)
```

## Benefits

- No network calls (faster, more resilient)
- No single point of failure
- Industry standard approach
