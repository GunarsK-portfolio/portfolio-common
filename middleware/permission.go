package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Permission levels (hierarchical)
var levelValues = map[string]int{
	"none":   0,
	"read":   1,
	"edit":   2,
	"delete": 3,
}

// HasPermission checks if user level meets required level
func HasPermission(userLevel, requiredLevel string) bool {
	return levelValues[userLevel] >= levelValues[requiredLevel]
}

// RequirePermission returns middleware that checks user has required permission
func RequirePermission(resource, level string) gin.HandlerFunc {
	return func(c *gin.Context) {
		scopes, exists := c.Get("scopes")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		scopesMap, ok := scopes.(map[string]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid scopes format"})
			return
		}

		userLevel := scopesMap[resource]
		if !HasPermission(userLevel, level) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":    "insufficient permissions",
				"resource": resource,
				"required": level,
				"have":     userLevel,
			})
			return
		}

		c.Next()
	}
}
