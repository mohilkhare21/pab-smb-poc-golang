package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mohilkhare21/pab-smb-poc-golang/backend/internal/auth"
	"github.com/mohilkhare21/pab-smb-poc-golang/backend/internal/database"
	"github.com/mohilkhare21/pab-smb-poc-golang/backend/internal/models"
)

// InvitationHandler handles invitation-related requests
type InvitationHandler struct {
	databaseProvider database.DatabaseProvider
	authProvider     auth.AuthProvider
}

// NewInvitationHandler creates a new invitation handler
func NewInvitationHandler(databaseProvider database.DatabaseProvider, authProvider auth.AuthProvider) *InvitationHandler {
	return &InvitationHandler{
		databaseProvider: databaseProvider,
		authProvider:     authProvider,
	}
}

// CreateInvitation handles creating user invitations
func (h *InvitationHandler) CreateInvitation(c *gin.Context) {
	var req models.InviteUserRequest
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

	// Check if user is admin
	if currentUser.Role != "admin" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "Only admins can invite users",
		})
		return
	}

	var createdInvitations []gin.H

	// Create invitations for each email
	for _, email := range req.Emails {
		// Check if user already exists
		existingUser, err := h.databaseProvider.GetUserByEmail(c.Request.Context(), email)
		if err == nil && existingUser != nil {
			// User already exists, skip
			continue
		}

		// Create invitation
		invitation := &database.Invitation{
			ID:        uuid.New().String(),
			Email:     email,
			CompanyID: currentUser.CompanyID,
			InvitedBy: currentUser.UserID,
			Token:     uuid.New().String(),
			Status:    "pending",
			ExpiresAt: time.Now().AddDate(0, 0, 7), // 7 days expiry
		}

		if err := h.databaseProvider.CreateInvitation(c.Request.Context(), invitation); err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Error:   "Failed to create invitation for " + email,
			})
			return
		}

		// Send invitation email
		if err := h.authProvider.SendInvitation(c.Request.Context(), email, currentUser.CompanyID, currentUser.UserID); err != nil {
			// Log error but don't fail the request
			// In a real implementation, you might want to handle this differently
		}

		createdInvitations = append(createdInvitations, gin.H{
			"id":        invitation.ID,
			"email":     invitation.Email,
			"status":    invitation.Status,
			"expires_at": invitation.ExpiresAt,
		})
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Invitations created successfully",
		Data: gin.H{
			"invitations": createdInvitations,
		},
	})
}

// GetInvitations handles getting all invitations for a company
func (h *InvitationHandler) GetInvitations(c *gin.Context) {
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

	// Check if user is admin
	if currentUser.Role != "admin" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "Only admins can view invitations",
		})
		return
	}

	// Get invitations for the company
	invitations, err := h.databaseProvider.GetInvitationsByCompany(c.Request.Context(), currentUser.CompanyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to get invitations",
		})
		return
	}

	// Convert to response format
	var invitationList []gin.H
	for _, invitation := range invitations {
		invitationList = append(invitationList, gin.H{
			"id":         invitation.ID,
			"email":      invitation.Email,
			"invited_by": invitation.InvitedBy,
			"status":     invitation.Status,
			"expires_at": invitation.ExpiresAt,
			"created_at": invitation.CreatedAt,
			"accepted_at": invitation.AcceptedAt,
		})
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: gin.H{
			"invitations": invitationList,
		},
	})
}

// DeleteInvitation handles deleting an invitation
func (h *InvitationHandler) DeleteInvitation(c *gin.Context) {
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
	invitationID := c.Param("id")

	// Check if user is admin
	if currentUser.Role != "admin" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "Only admins can delete invitations",
		})
		return
	}

	// Get invitation
	invitation, err := h.databaseProvider.GetInvitation(c.Request.Context(), invitationID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Invitation not found",
		})
		return
	}

	// Check if invitation belongs to the company
	if invitation.CompanyID != currentUser.CompanyID {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "Access denied",
		})
		return
	}

	// Delete invitation
	if err := h.databaseProvider.DeleteInvitation(c.Request.Context(), invitationID); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to delete invitation",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Invitation deleted successfully",
	})
}

// AcceptInvitation handles accepting an invitation
func (h *InvitationHandler) AcceptInvitation(c *gin.Context) {
	token := c.Param("token")

	// Get invitation by token
	invitation, err := h.databaseProvider.GetInvitationByToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Invalid invitation token",
		})
		return
	}

	// Check if invitation is expired
	if time.Now().After(invitation.ExpiresAt) {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invitation has expired",
		})
		return
	}

	// Check if invitation is already accepted
	if invitation.Status == "accepted" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invitation has already been accepted",
		})
		return
	}

	// Get user from context (user should be authenticated to accept invitation)
	userContext, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "User not authenticated",
		})
		return
	}

	currentUser := userContext.(models.UserContext)

	// Check if user email matches invitation email
	if currentUser.Email != invitation.Email {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "Email does not match invitation",
		})
		return
	}

	// Update user with company ID
	user, err := h.databaseProvider.GetUser(c.Request.Context(), currentUser.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to get user data",
		})
		return
	}

	user.CompanyID = invitation.CompanyID
	user.Role = "user" // Default role for invited users
	if err := h.databaseProvider.UpdateUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to update user",
		})
		return
	}

	// Update invitation status
	invitation.Status = "accepted"
	invitation.AcceptedAt = time.Now()
	if err := h.databaseProvider.UpdateInvitation(c.Request.Context(), invitation); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to update invitation",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Invitation accepted successfully",
		Data: gin.H{
			"company_id": invitation.CompanyID,
		},
	})
}

