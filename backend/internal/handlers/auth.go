package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mohilkhare21/pab-smb-poc-golang/backend/internal/auth"
	"github.com/mohilkhare21/pab-smb-poc-golang/backend/internal/database"
	"github.com/mohilkhare21/pab-smb-poc-golang/backend/internal/models"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	authProvider     auth.AuthProvider
	databaseProvider database.DatabaseProvider
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authProvider auth.AuthProvider, databaseProvider database.DatabaseProvider) *AuthHandler {
	return &AuthHandler{
		authProvider:     authProvider,
		databaseProvider: databaseProvider,
	}
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request data: " + err.Error(),
		})
		return
	}

	// Authenticate user
	user, err := h.authProvider.Authenticate(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Invalid credentials",
		})
		return
	}

	// Get user from database to get company info
	dbUser, err := h.databaseProvider.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to get user data",
		})
		return
	}

	// Generate JWT token
	token, err := h.authProvider.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Login successful",
		Data: gin.H{
			"token": token,
			"user": gin.H{
				"id":         user.ID,
				"email":      user.Email,
				"name":       user.Name,
				"picture":    user.Picture,
				"company_id": dbUser.CompanyID,
				"role":       dbUser.Role,
			},
		},
	})
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request data: " + err.Error(),
		})
		return
	}

	// Check if user already exists
	existingUser, err := h.databaseProvider.GetUserByEmail(c.Request.Context(), req.Email)
	if err == nil && existingUser != nil {
		c.JSON(http.StatusConflict, models.APIResponse{
			Success: false,
			Error:   "User already exists",
		})
		return
	}

	// Register user with auth provider
	user, err := h.authProvider.Register(c.Request.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to register user: " + err.Error(),
		})
		return
	}

	// Create user in database
	dbUser := &database.User{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Picture:   user.Picture,
		Role:      string(auth.RoleAdmin), // First user becomes admin
		IsActive:  true,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if err := h.databaseProvider.CreateUser(c.Request.Context(), dbUser); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to create user in database",
		})
		return
	}

	// Generate JWT token
	token, err := h.authProvider.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Registration successful",
		Data: gin.H{
			"token": token,
			"user": gin.H{
				"id":      user.ID,
				"email":   user.Email,
				"name":    user.Name,
				"picture": user.Picture,
				"role":    auth.RoleAdmin,
			},
		},
	})
}

// AuthenticateWithToken handles authentication with JWT token
func (h *AuthHandler) AuthenticateWithToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Authorization header required",
		})
		return
	}

	// Extract token from "Bearer <token>"
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Invalid authorization header format",
		})
		return
	}

	token := tokenParts[1]

	// Validate token
	user, err := h.authProvider.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Invalid token",
		})
		return
	}

	// Get user from database
	dbUser, err := h.databaseProvider.GetUser(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Authentication successful",
		Data: gin.H{
			"user": gin.H{
				"id":         dbUser.ID,
				"email":      dbUser.Email,
				"name":       dbUser.Name,
				"picture":    dbUser.Picture,
				"company_id": dbUser.CompanyID,
				"role":       dbUser.Role,
				"is_active":  dbUser.IsActive,
			},
		},
	})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request data",
		})
		return
	}

	// Refresh token
	newToken, err := h.authProvider.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Invalid refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Token refreshed successfully",
		Data: gin.H{
			"token": newToken,
		},
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// In a stateless JWT system, logout is typically handled on the client side
	// by removing the token. However, you could implement a blacklist for tokens.
	
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Logout successful",
	})
}

// ResetPassword handles password reset request
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid email address",
		})
		return
	}

	// Send password reset email
	err := h.authProvider.ResetPassword(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to send password reset email",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Password reset email sent",
	})
}

// ChangePassword handles password change
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request data",
		})
		return
	}

	// Get user from context (set by auth middleware)
	userContext, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "User not authenticated",
		})
		return
	}

	user := userContext.(models.UserContext)

	// Change password
	err := h.authProvider.ChangePassword(c.Request.Context(), user.UserID, req.OldPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Failed to change password",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}

