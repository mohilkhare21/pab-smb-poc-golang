package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mohilkhare21/pab-smb-poc-golang/backend/internal/database"
	"github.com/mohilkhare21/pab-smb-poc-golang/backend/internal/models"
)

// BrowserShortcutHandler handles browser shortcut-related requests
type BrowserShortcutHandler struct {
	databaseProvider database.DatabaseProvider
}

// NewBrowserShortcutHandler creates a new browser shortcut handler
func NewBrowserShortcutHandler(databaseProvider database.DatabaseProvider) *BrowserShortcutHandler {
	return &BrowserShortcutHandler{
		databaseProvider: databaseProvider,
	}
}

// GetShortcuts handles getting all browser shortcuts for a company
func (h *BrowserShortcutHandler) GetShortcuts(c *gin.Context) {
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

	// Get shortcuts for the company
	shortcuts, err := h.databaseProvider.GetBrowserShortcutsByCompany(c.Request.Context(), user.CompanyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to get shortcuts",
		})
		return
	}

	// Convert to response format
	var shortcutList []gin.H
	for _, shortcut := range shortcuts {
		shortcutList = append(shortcutList, gin.H{
			"id":          shortcut.ID,
			"name":        shortcut.Name,
			"url":         shortcut.URL,
			"icon":        shortcut.Icon,
			"description": shortcut.Description,
			"order":       shortcut.Order,
			"is_active":   shortcut.IsActive,
		})
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: gin.H{
			"shortcuts": shortcutList,
		},
	})
}

// CreateShortcut handles creating a new browser shortcut
func (h *BrowserShortcutHandler) CreateShortcut(c *gin.Context) {
	var req models.BrowserShortcutRequest
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

	user := userContext.(models.UserContext)

	// Check if user is admin
	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "Only admins can create shortcuts",
		})
		return
	}

	// Create shortcut
	shortcut := &database.BrowserShortcut{
		ID:          uuid.New().String(),
		CompanyID:   user.CompanyID,
		Name:        req.Name,
		URL:         req.URL,
		Icon:        req.Icon,
		Description: req.Description,
		Order:       req.Order,
		IsActive:    true,
	}

	if err := h.databaseProvider.CreateBrowserShortcut(c.Request.Context(), shortcut); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to create shortcut",
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Shortcut created successfully",
		Data: gin.H{
			"shortcut": gin.H{
				"id":          shortcut.ID,
				"name":        shortcut.Name,
				"url":         shortcut.URL,
				"icon":        shortcut.Icon,
				"description": shortcut.Description,
				"order":       shortcut.Order,
				"is_active":   shortcut.IsActive,
			},
		},
	})
}

// UpdateShortcut handles updating a browser shortcut
func (h *BrowserShortcutHandler) UpdateShortcut(c *gin.Context) {
	var req models.BrowserShortcutRequest
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

	user := userContext.(models.UserContext)
	shortcutID := c.Param("id")

	// Check if user is admin
	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "Only admins can update shortcuts",
		})
		return
	}

	// Get shortcut
	shortcut, err := h.databaseProvider.GetBrowserShortcut(c.Request.Context(), shortcutID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Shortcut not found",
		})
		return
	}

	// Check if shortcut belongs to the company
	if shortcut.CompanyID != user.CompanyID {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "Access denied",
		})
		return
	}

	// Update shortcut fields
	if req.Name != "" {
		shortcut.Name = req.Name
	}
	if req.URL != "" {
		shortcut.URL = req.URL
	}
	if req.Icon != "" {
		shortcut.Icon = req.Icon
	}
	if req.Description != "" {
		shortcut.Description = req.Description
	}
	if req.Order != 0 {
		shortcut.Order = req.Order
	}

	// Save updated shortcut
	if err := h.databaseProvider.UpdateBrowserShortcut(c.Request.Context(), shortcut); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to update shortcut",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Shortcut updated successfully",
		Data: gin.H{
			"shortcut": gin.H{
				"id":          shortcut.ID,
				"name":        shortcut.Name,
				"url":         shortcut.URL,
				"icon":        shortcut.Icon,
				"description": shortcut.Description,
				"order":       shortcut.Order,
				"is_active":   shortcut.IsActive,
			},
		},
	})
}

// DeleteShortcut handles deleting a browser shortcut
func (h *BrowserShortcutHandler) DeleteShortcut(c *gin.Context) {
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
	shortcutID := c.Param("id")

	// Check if user is admin
	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "Only admins can delete shortcuts",
		})
		return
	}

	// Get shortcut
	shortcut, err := h.databaseProvider.GetBrowserShortcut(c.Request.Context(), shortcutID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Shortcut not found",
		})
		return
	}

	// Check if shortcut belongs to the company
	if shortcut.CompanyID != user.CompanyID {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "Access denied",
		})
		return
	}

	// Delete shortcut
	if err := h.databaseProvider.DeleteBrowserShortcut(c.Request.Context(), shortcutID); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to delete shortcut",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Shortcut deleted successfully",
	})
}

