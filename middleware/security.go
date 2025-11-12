package middleware

import (
	"slices"

	"github.com/gin-gonic/gin"
)

// SecurityMiddleware provides CORS validation and security headers
type SecurityMiddleware struct {
	allowedOrigins   []string
	allowedMethods   string
	allowedHeaders   string
	allowCredentials bool
}

// NewSecurityMiddleware creates a new security middleware with CORS configuration
func NewSecurityMiddleware(allowedOrigins []string, allowedMethods string, allowedHeaders string, allowCredentials bool) *SecurityMiddleware {
	if len(allowedOrigins) == 0 {
		panic("allowedOrigins must contain at least one origin")
	}
	for _, origin := range allowedOrigins {
		if origin == "" {
			panic("allowedOrigins must not contain empty strings")
		}
	}

	return &SecurityMiddleware{
		allowedOrigins:   allowedOrigins,
		allowedMethods:   allowedMethods,
		allowedHeaders:   allowedHeaders,
		allowCredentials: allowCredentials,
	}
}

// Apply returns a Gin middleware that adds security headers and validates CORS
func (m *SecurityMiddleware) Apply() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// CORS validation - only set headers if origin is allowed
		allowed := slices.Contains(m.allowedOrigins, origin)

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Methods", m.allowedMethods)
			c.Writer.Header().Set("Access-Control-Allow-Headers", m.allowedHeaders)
			c.Writer.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
			if m.allowCredentials {
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			}
		}

		// Standard security headers (applied to all requests)
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")

		// Handle preflight OPTIONS requests
		if c.Request.Method == "OPTIONS" {
			if allowed {
				c.AbortWithStatus(204)
			} else {
				c.AbortWithStatus(403)
			}
			return
		}

		c.Next()
	}
}
