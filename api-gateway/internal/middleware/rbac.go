package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole ensures user has one of the allowed roles
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool)
	for _, r := range allowedRoles {
		allowed[r] = true
	}

	return func(c *gin.Context) {
		role := GetRole(c)
		if role == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "role not found"})
			return
		}

		if !allowed[role] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "insufficient privileges",
				"role":    role,
				"allowed": allowedRoles,
			})
			return
		}

		c.Next()
	}
}

// Shortcut for admin only
func AdminOnly() gin.HandlerFunc {
	return RequireRole("admin")
}

// Shortcut for any authenticated user
func AnyUser() gin.HandlerFunc {
	return RequireRole("customer", "admin", "vendor")
}
