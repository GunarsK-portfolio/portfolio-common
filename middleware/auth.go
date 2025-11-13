package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware provides JWT token validation via auth-service
type AuthMiddleware struct {
	authServiceURL string
	timeout        time.Duration
	client         *http.Client
}

// AuthMiddlewareOption is a functional option for configuring AuthMiddleware
type AuthMiddlewareOption func(*AuthMiddleware)

// WithTimeout sets a custom timeout for auth service requests.
// The timeout must be positive. Zero or negative values will be ignored
// and the default timeout (5s) will be used instead.
func WithTimeout(timeout time.Duration) AuthMiddlewareOption {
	return func(m *AuthMiddleware) {
		if timeout > 0 {
			m.timeout = timeout
		}
	}
}

// NewAuthMiddleware creates a new auth middleware instance with optional configuration
// Default timeout is 5 seconds if not specified
// Usage: NewAuthMiddleware(url) or NewAuthMiddleware(url, WithTimeout(10*time.Second))
func NewAuthMiddleware(authServiceURL string, opts ...AuthMiddlewareOption) *AuthMiddleware {
	m := &AuthMiddleware{
		authServiceURL: authServiceURL,
		timeout:        5 * time.Second, // Default timeout
	}

	// Apply optional configurations
	for _, opt := range opts {
		opt(m)
	}

	// Initialize HTTP client with configured timeout for connection pooling
	m.client = &http.Client{
		Timeout: m.timeout,
	}

	return m
}

// ValidateToken returns a Gin middleware that validates JWT tokens via auth-service
func (m *AuthMiddleware) ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized - no token provided"})
			c.Abort()
			return
		}

		// Validate token with auth service and get TTL + claims
		ttl, claims, err := m.validateWithAuthService(token)
		if err != nil {
			slog.Error("token validation failed",
				"error", err,
				"auth_service_url", m.authServiceURL,
				"path", c.Request.URL.Path,
			)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized - invalid token"})
			c.Abort()
			return
		}

		if ttl <= 0 {
			slog.Warn("token validation returned non-positive TTL",
				"ttl", ttl,
				"path", c.Request.URL.Path,
			)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized - invalid token"})
			c.Abort()
			return
		}

		// Store TTL and user claims in context for downstream handlers
		c.Set("token_ttl", ttl)
		if claims != nil {
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
		}

		c.Next()
	}
}

// TokenClaims represents JWT claims extracted from validation response
type TokenClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}

func (m *AuthMiddleware) validateWithAuthService(token string) (int64, *TokenClaims, error) {
	validateURL := fmt.Sprintf("%s/auth/validate", m.authServiceURL)

	reqBody, _ := json.Marshal(map[string]string{"token": token})

	req, err := http.NewRequest(http.MethodPost, validateURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create auth request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("auth service request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Limit response size to 1KB to prevent DoS attacks
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024))
	if err != nil {
		return 0, nil, fmt.Errorf("failed to read auth response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		// Include limited response body for debugging (max 512 bytes)
		bodyPreview := string(body)
		if len(bodyPreview) > 512 {
			bodyPreview = bodyPreview[:512] + "..."
		}
		return 0, nil, fmt.Errorf("auth service returned status %d, body: %s", resp.StatusCode, bodyPreview)
	}

	var result struct {
		Valid      bool   `json:"valid"`
		TTLSeconds int64  `json:"ttl_seconds"`
		UserID     int64  `json:"user_id"`
		Username   string `json:"username"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, nil, fmt.Errorf("failed to parse auth response: %w", err)
	}

	if !result.Valid {
		return 0, nil, fmt.Errorf("token validation failed")
	}

	claims := &TokenClaims{
		UserID:   result.UserID,
		Username: result.Username,
	}

	return result.TTLSeconds, claims, nil
}

// AddTTLHeader returns a middleware that adds X-Token-TTL header to responses
// This should be added after ValidateToken middleware
func (m *AuthMiddleware) AddTTLHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// After request processing, add TTL header if available in context
		if ttl, exists := c.Get("token_ttl"); exists {
			if ttlValue, ok := ttl.(int64); ok && ttlValue > 0 {
				c.Header("X-Token-TTL", fmt.Sprintf("%d", ttlValue))
			}
		}
	}
}

func extractToken(c *gin.Context) string {
	bearerToken := c.GetHeader("Authorization")
	parts := strings.Split(bearerToken, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}
	return ""
}
