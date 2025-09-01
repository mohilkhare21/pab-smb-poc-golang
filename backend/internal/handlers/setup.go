package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohilkhare21/pab-smb-poc-golang/backend/internal/database"
	"github.com/mohilkhare21/pab-smb-poc-golang/backend/internal/models"
)

// SetupHandler handles company setup and configuration
type SetupHandler struct {
	databaseProvider database.DatabaseProvider
}

// NewSetupHandler creates a new setup handler
func NewSetupHandler(db database.DatabaseProvider) *SetupHandler {
	return &SetupHandler{
		databaseProvider: db,
	}
}

// GetSetupProgress returns the setup progress for a company
func (h *SetupHandler) GetSetupProgress(c *gin.Context) {
	userContext, exists := c.Get("user_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "User context not found",
		})
		return
	}

	uc := userContext.(models.UserContext)
	
	progress, err := h.databaseProvider.GetSetupProgress(c.Request.Context(), uc.CompanyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to get setup progress",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    progress,
	})
}

// UpdateSetupStep updates the setup step for a company
func (h *SetupHandler) UpdateSetupStep(c *gin.Context) {
	userContext, exists := c.Get("user_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "User context not found",
		})
		return
	}

	uc := userContext.(models.UserContext)

	var req struct {
		Step     string `json:"step" binding:"required"`
		Progress int    `json:"progress" binding:"required,min=0,max=100"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request data",
		})
		return
	}

	err := h.databaseProvider.UpdateSetupStep(c.Request.Context(), uc.CompanyID, req.Step, req.Progress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to update setup step",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Setup step updated successfully",
	})
}

// GetCompanyStats returns company statistics
func (h *SetupHandler) GetCompanyStats(c *gin.Context) {
	userContext, exists := c.Get("user_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "User context not found",
		})
		return
	}

	uc := userContext.(models.UserContext)

	// Get company
	company, err := h.databaseProvider.GetCompany(c.Request.Context(), uc.CompanyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to get company",
		})
		return
	}

	// Get user counts
	totalUsers, err := h.databaseProvider.CountUsersByCompany(c.Request.Context(), uc.CompanyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to count users",
		})
		return
	}

	activeUsers, err := h.databaseProvider.CountActiveUsersByCompany(c.Request.Context(), uc.CompanyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to count active users",
		})
		return
	}

	invitedUsers, err := h.databaseProvider.CountInvitedUsersByCompany(c.Request.Context(), uc.CompanyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to count invited users",
		})
		return
	}

	pendingInvitations, err := h.databaseProvider.CountPendingInvitationsByCompany(c.Request.Context(), uc.CompanyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to count pending invitations",
		})
		return
	}

	// Get subscription stats
	subscription, err := h.databaseProvider.GetSubscriptionByCompany(c.Request.Context(), uc.CompanyID)
	if err != nil {
		// Subscription might not exist yet
		subscription = &database.Subscription{
			MaxUsers: 20, // Default max users
		}
	}

	// Get configuration status
	configStatus, err := h.databaseProvider.GetCompanyConfigurationStatus(c.Request.Context(), uc.CompanyID)
	if err != nil {
		configStatus = make(map[string]bool)
	}

	// Calculate setup progress
	setupProgress := 0
	if company.Domain != "" {
		setupProgress += 20
	}
	if company.ColorTheme != "" {
		setupProgress += 20
	}
	if company.UsersInvited {
		setupProgress += 20
	}
	if company.SubscriptionActive {
		setupProgress += 20
	}
	if company.DownloadReady {
		setupProgress += 20
	}

	stats := models.CompanyStatsResponse{
		TotalUsers:        totalUsers,
		ActiveUsers:       activeUsers,
		InvitedUsers:      invitedUsers,
		PendingInvitations: pendingInvitations,
		MaxUsers:          subscription.MaxUsers,
		SetupProgress:     setupProgress,
		ConfigurationStatus: configStatus,
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    stats,
	})
}

// UpdateConfigurationStatus updates configuration status for a feature
func (h *SetupHandler) UpdateConfigurationStatus(c *gin.Context) {
	userContext, exists := c.Get("user_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "User context not found",
		})
		return
	}

	uc := userContext.(models.UserContext)

	var req models.ConfigurationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request data",
		})
		return
	}

	err := h.databaseProvider.UpdateCompanyConfigurationStatus(c.Request.Context(), uc.CompanyID, req.Feature, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to update configuration status",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Configuration status updated successfully",
	})
}

// GenerateShortcuts generates suggested shortcuts for a company domain
func (h *SetupHandler) GenerateShortcuts(c *gin.Context) {
	userContext, exists := c.Get("user_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "User context not found",
		})
		return
	}

	uc := userContext.(models.UserContext)

	var req models.GenerateShortcutsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request data",
		})
		return
	}

	err := h.databaseProvider.GenerateShortcutsForDomain(c.Request.Context(), uc.CompanyID, req.Domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to generate shortcuts",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Shortcuts generated successfully",
	})
}

// NudgeUsers sends reminders to invited users
func (h *SetupHandler) NudgeUsers(c *gin.Context) {
	userContext, exists := c.Get("user_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "User context not found",
		})
		return
	}

	uc := userContext.(models.UserContext)

	var req models.NudgeUsersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request data",
		})
		return
	}

	// Get pending invitations for the company
	invitations, err := h.databaseProvider.GetPendingInvitationsByCompany(c.Request.Context(), uc.CompanyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to get pending invitations",
		})
		return
	}

	// Resend invitations
	for _, invitation := range invitations {
		err := h.databaseProvider.ResendInvitation(c.Request.Context(), invitation.ID)
		if err != nil {
			// Log error but continue with other invitations
			continue
		}
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User nudges sent successfully",
	})
}

// GetDownloadInfo returns download information for the custom browser
func (h *SetupHandler) GetDownloadInfo(c *gin.Context) {
	userContext, exists := c.Get("user_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "User context not found",
		})
		return
	}

	uc := userContext.(models.UserContext)

	// Check if company setup is complete
	company, err := h.databaseProvider.GetCompany(c.Request.Context(), uc.CompanyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to get company",
		})
		return
	}

	if !company.DownloadReady {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Company setup not complete. Download not ready.",
		})
		return
	}

	downloadInfo := models.DownloadInfoResponse{
		DownloadURL: "https://download.pab-smb.com/browser/latest",
		Version:     "1.0.0",
		ReleaseDate: time.Now().Format("2006-01-02"),
		FileSize:    "45.2 MB",
		SupportedOS: []string{"macOS", "Windows", "Linux"},
		InstallationInstructions: "Download and run the installer. Follow the setup wizard to complete installation.",
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    downloadInfo,
	})
}

