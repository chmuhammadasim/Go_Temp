package middleware

import (
	"net/http"
	"strings"

	"go-backend/internal/models"
	"go-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware(jwtService *utils.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_username", claims.Username)
		c.Set("user_role", claims.Role)
		c.Set("claims", claims)

		c.Next()
	}
}

// RequireRole middleware checks if user has required role
func RequireRole(requiredRoles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User role not found in context",
			})
			c.Abort()
			return
		}

		role, ok := userRole.(models.Role)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid user role type",
			})
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		hasPermission := false
		for _, requiredRole := range requiredRoles {
			if hasRolePermission(role, requiredRole) {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin middleware checks if user has admin role
func RequireAdmin() gin.HandlerFunc {
	return RequireRole(models.RoleAdmin)
}

// RequireModerator middleware checks if user has moderator or admin role
func RequireModerator() gin.HandlerFunc {
	return RequireRole(models.RoleAdmin, models.RoleModerator)
}

// RequireOwnerOrAdmin middleware checks if user is the owner of the resource or admin
func RequireOwnerOrAdmin(getUserIDFunc func(*gin.Context) uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User ID not found in context",
			})
			c.Abort()
			return
		}

		currentUserID, ok := userID.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid user ID type",
			})
			c.Abort()
			return
		}

		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User role not found in context",
			})
			c.Abort()
			return
		}

		role, ok := userRole.(models.Role)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid user role type",
			})
			c.Abort()
			return
		}

		// Admin can access everything
		if role == models.RoleAdmin {
			c.Next()
			return
		}

		// Check if user is the owner
		resourceUserID := getUserIDFunc(c)
		if currentUserID == resourceUserID {
			c.Next()
			return
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": "Access denied: you can only access your own resources",
		})
		c.Abort()
	}
}

// hasRolePermission checks if a user role has permission for a required role
func hasRolePermission(userRole, requiredRole models.Role) bool {
	switch requiredRole {
	case models.RoleAdmin:
		return userRole == models.RoleAdmin
	case models.RoleModerator:
		return userRole == models.RoleAdmin || userRole == models.RoleModerator
	case models.RoleUser:
		return true // All authenticated users have user permission
	default:
		return false
	}
}
