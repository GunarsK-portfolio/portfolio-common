package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware provides JWT token validation via auth-service
type AuthMiddleware struct {
	authServiceURL string
}

// NewAuthMiddleware creates a new auth middleware instance
func NewAuthMiddleware(authServiceURL string) *AuthMiddleware {
	return &AuthMiddleware{authServiceURL: authServiceURL}
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

		// Validate token with auth service and get TTL
		ttl, err := m.validateWithAuthService(token)
		if err != nil || ttl <= 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized - invalid token"})
			c.Abort()
			return
		}

		// Store TTL in context for response middleware
		c.Set("token_ttl", ttl)

		c.Next()
	}
}

func (m *AuthMiddleware) validateWithAuthService(token string) (int64, error) {
	validateURL := fmt.Sprintf("%s/auth/validate", m.authServiceURL)

	reqBody, _ := json.Marshal(map[string]string{"token": token})

	req, err := http.NewRequest(http.MethodPost, validateURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return 0, nil
	}

	var result struct {
		Valid      bool  `json:"valid"`
		TTLSeconds int64 `json:"ttl_seconds"`
	}
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	if !result.Valid {
		return 0, nil
	}

	return result.TTLSeconds, nil
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
