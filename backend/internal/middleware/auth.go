package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mohilkhare21/pab-smb-poc-golang/backend/internal/auth"
	"github.com/mohilkhare21/pab-smb-poc-golang/backend/internal/database"
	"github.com/mohilkhare21/pab-smb-poc-golang/backend/internal/models"
)

// AuthMiddleware handles authentication middleware
type AuthMiddleware struct {
	authProvider     auth.AuthProvider
	databaseProvider database.DatabaseProvider
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(authProvider auth.AuthProvider) *AuthMiddleware {
	return &AuthMiddleware{
		authProvider: authProvider,
	}
}

// SetDatabaseProvider sets the database provider for the middleware
func (m *AuthMiddleware) SetDatabaseProvider(dbProvider database.DatabaseProvider) {
	m.databaseProvider = dbProvider
}

// Authenticate middleware validates JWT tokens and sets user context
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Error:   "Authorization header required",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Error:   "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// Validate token
		user, err := m.authProvider.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Error:   "Invalid token",
			})
			c.Abort()
			return
		}

		// Get user from database to get complete user info
		if m.databaseProvider != nil {
			dbUser, err := m.databaseProvider.GetUser(c.Request.Context(), user.ID)
			if err != nil {
				c.JSON(http.StatusUnauthorized, models.APIResponse{
					Success: false,
					Error:   "User not found",
				})
				c.Abort()
				return
			}

			// Check if user is active
			if !dbUser.IsActive {
				c.JSON(http.StatusForbidden, models.APIResponse{
					Success: false,
					Error:   "User account is inactive",
				})
				c.Abort()
				return
			}

			// Set user context
			userContext := models.UserContext{
				UserID:    dbUser.ID,
				Email:     dbUser.Email,
				CompanyID: dbUser.CompanyID,
				Role:      dbUser.Role,
			}
			c.Set("user", userContext)
		} else {
			// Fallback if database provider is not set
			userContext := models.UserContext{
				UserID: user.ID,
				Email:  user.Email,
				Role:   "user", // Default role
			}
			c.Set("user", userContext)
		}

		c.Next()
	}
}

// RequireRole middleware checks if user has required role
func (m *AuthMiddleware) RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context (set by Authenticate middleware)
		userContext, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Error:   "User not authenticated",
			})
			c.Abort()
			return
		}

		user := userContext.(models.UserContext)

		// Check if user has required role
		if user.Role != requiredRole {
			c.JSON(http.StatusForbidden, models.APIResponse{
				Success: false,
				Error:   "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyRole middleware checks if user has any of the required roles
func (m *AuthMiddleware) RequireAnyRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context (set by Authenticate middleware)
		userContext, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Error:   "User not authenticated",
			})
			c.Abort()
			return
		}

		user := userContext.(models.UserContext)

		// Check if user has any of the required roles
		hasRole := false
		for _, role := range requiredRoles {
			if user.Role == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, models.APIResponse{
				Success: false,
				Error:   "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireCompanyAccess middleware ensures user has access to the specified company
func (m *AuthMiddleware) RequireCompanyAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context (set by Authenticate middleware)
		userContext, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Error:   "User not authenticated",
			})
			c.Abort()
			return
		}

		user := userContext.(models.UserContext)

		// Check if user has a company ID
		if user.CompanyID == "" {
			c.JSON(http.StatusForbidden, models.APIResponse{
				Success: false,
				Error:   "User not associated with any company",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuth middleware validates JWT tokens if present but doesn't require them
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No auth header, continue without user context
			c.Next()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			// Invalid format, continue without user context
			c.Next()
			return
		}

		token := tokenParts[1]

		// Validate token
		user, err := m.authProvider.ValidateToken(token)
		if err != nil {
			// Invalid token, continue without user context
			c.Next()
			return
		}

		// Get user from database to get complete user info
		if m.databaseProvider != nil {
			dbUser, err := m.databaseProvider.GetUser(c.Request.Context(), user.ID)
			if err == nil && dbUser.IsActive {
				// Set user context
				userContext := models.UserContext{
					UserID:    dbUser.ID,
					Email:     dbUser.Email,
					CompanyID: dbUser.CompanyID,
					Role:      dbUser.Role,
				}
				c.Set("user", userContext)
			}
		} else {
			// Fallback if database provider is not set
			userContext := models.UserContext{
				UserID: user.ID,
				Email:  user.Email,
				Role:   "user", // Default role
			}
			c.Set("user", userContext)
		}

		c.Next()
	}
}

