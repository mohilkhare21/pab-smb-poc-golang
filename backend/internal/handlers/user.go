package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohilkhare21/pab-smb-poc-golang/backend/internal/auth"
	"github.com/mohilkhare21/pab-smb-poc-golang/backend/internal/database"
	"github.com/mohilkhare21/pab-smb-poc-golang/backend/internal/models"
)

// UserHandler handles user-related requests
type UserHandler struct {
	databaseProvider database.DatabaseProvider
	authProvider     auth.AuthProvider
}

// NewUserHandler creates a new user handler
func NewUserHandler(databaseProvider database.DatabaseProvider, authProvider auth.AuthProvider) *UserHandler {
	return &UserHandler{
		databaseProvider: databaseProvider,
		authProvider:     authProvider,
	}
}

// GetUsers handles getting all users for a company
func (h *UserHandler) GetUsers(c *gin.Context) {
	// Get user from context
	userContext, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "User not authenticated",
		})
		return
	}

	user := userContext.(models.UserContext)

	// Get users for the company
	users, err := h.databaseProvider.GetUsersByCompany(c.Request.Context(), user.CompanyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to get users",
		})
		return
	}

	// Convert to response format
	var userList []gin.H
	for _, u := range users {
		userList = append(userList, gin.H{
			"id":           u.ID,
			"email":        u.Email,
			"name":         u.Name,
			"picture":      u.Picture,
			"role":         u.Role,
			"is_active":    u.IsActive,
			"created_at":   u.CreatedAt,
			"updated_at":   u.UpdatedAt,
			"last_login_at": u.LastLoginAt,
			"onboarded":    u.Onboarded,
		})
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: gin.H{
			"users": userList,
		},
	})
}

// GetUser handles getting a specific user
func (h *UserHandler) GetUser(c *gin.Context) {
	// Get user from context
	userContext, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "User not authenticated",
		})
		return
	}

	currentUser := userContext.(models.UserContext)
	userID := c.Param("id")

	// Get user from database
	user, err := h.databaseProvider.GetUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	// Check if user belongs to the same company
	if user.CompanyID != currentUser.CompanyID {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "Access denied",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: gin.H{
			"user": gin.H{
				"id":           user.ID,
				"email":        user.Email,
				"name":         user.Name,
				"picture":      user.Picture,
				"role":         user.Role,
				"is_active":    user.IsActive,
				"created_at":   user.CreatedAt,
				"updated_at":   user.UpdatedAt,
				"last_login_at": user.LastLoginAt,
				"onboarded":    user.Onboarded,
			},
		},
	})
}

// UpdateUser handles updating a user
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request data: " + err.Error(),
		})
		return
	}

	// Get user from context
	userContext, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "User not authenticated",
		})
		return
	}

	currentUser := userContext.(models.UserContext)
	userID := c.Param("id")

	// Get user from database
	user, err := h.databaseProvider.GetUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	// Check if user belongs to the same company
	if user.CompanyID != currentUser.CompanyID {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "Access denied",
		})
		return
	}

	// Check if current user is admin or updating their own profile
	if currentUser.Role != "admin" && currentUser.UserID != userID {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "Only admins can update other users",
		})
		return
	}

	// Update user fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Role != "" {
		// Only admins can change roles
		if currentUser.Role != "admin" {
			c.JSON(http.StatusForbidden, models.APIResponse{
				Success: false,
				Error:   "Only admins can change user roles",
			})
			return
		}
		user.Role = string(req.Role)
	}
	if req.IsActive != nil {
		// Only admins can deactivate users
		if currentUser.Role != "admin" {
			c.JSON(http.StatusForbidden, models.APIResponse{
				Success: false,
				Error:   "Only admins can deactivate users",
			})
			return
		}
		user.IsActive = *req.IsActive
	}

	// Save updated user
	if err := h.databaseProvider.UpdateUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User updated successfully",
		Data: gin.H{
			"user": gin.H{
				"id":           user.ID,
				"email":        user.Email,
				"name":         user.Name,
				"picture":      user.Picture,
				"role":         user.Role,
				"is_active":    user.IsActive,
				"updated_at":   user.UpdatedAt,
			},
		},
	})
}

// DeleteUser handles deleting a user
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// Get user from context
	userContext, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "User not authenticated",
		})
		return
	}

	currentUser := userContext.(models.UserContext)
	userID := c.Param("id")

	// Check if current user is admin
	if currentUser.Role != "admin" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "Only admins can delete users",
		})
		return
	}

	// Get user from database
	user, err := h.databaseProvider.GetUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	// Check if user belongs to the same company
	if user.CompanyID != currentUser.CompanyID {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "Access denied",
		})
		return
	}

	// Prevent deleting the last admin
	if user.Role == "admin" {
		users, err := h.databaseProvider.GetUsersByCompany(c.Request.Context(), currentUser.CompanyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Error:   "Failed to check user count",
			})
			return
		}

		adminCount := 0
		for _, u := range users {
			if u.Role == "admin" && u.IsActive {
				adminCount++
			}
		}

		if adminCount <= 1 {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error:   "Cannot delete the last admin user",
			})
			return
		}
	}

	// Delete user from auth provider
	if err := h.authProvider.DeleteUser(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to delete user from auth provider",
		})
		return
	}

	// Delete user from database
	if err := h.databaseProvider.DeleteUser(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to delete user from database",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User deleted successfully",
	})
}
