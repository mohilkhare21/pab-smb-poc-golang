package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mohilkhare21/pab-smb-poc-golang/backend/internal/database"
	"github.com/mohilkhare21/pab-smb-poc-golang/backend/internal/models"
)

// CompanyHandler handles company-related requests
type CompanyHandler struct {
	databaseProvider database.DatabaseProvider
}

// NewCompanyHandler creates a new company handler
func NewCompanyHandler(databaseProvider database.DatabaseProvider) *CompanyHandler {
	return &CompanyHandler{
		databaseProvider: databaseProvider,
	}
}

// CreateCompany handles company creation
func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	var req models.CompanyCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request data: " + err.Error(),
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

	// Check if domain is already taken
	existingCompany, err := h.databaseProvider.GetCompanyByDomain(c.Request.Context(), req.Domain)
	if err == nil && existingCompany != nil {
		c.JSON(http.StatusConflict, models.APIResponse{
			Success: false,
			Error:   "Domain already taken",
		})
		return
	}

	// Create company
	company := &database.Company{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Domain:      req.Domain,
		ColorTheme:  req.ColorTheme,
		AdminUserID: user.UserID,
		Status:      "trial",
	}

	if err := h.databaseProvider.CreateCompany(c.Request.Context(), company); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to create company",
		})
		return
	}

	// Update user with company ID
	dbUser, err := h.databaseProvider.GetUser(c.Request.Context(), user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to get user data",
		})
		return
	}

	dbUser.CompanyID = company.ID
	if err := h.databaseProvider.UpdateUser(c.Request.Context(), dbUser); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Company created successfully",
		Data: gin.H{
			"company": gin.H{
				"id":          company.ID,
				"name":        company.Name,
				"domain":      company.Domain,
				"color_theme": company.ColorTheme,
				"status":      company.Status,
				"admin_user_id": company.AdminUserID,
			},
		},
	})
}

// GetCompany handles getting company details
func (h *CompanyHandler) GetCompany(c *gin.Context) {
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

	// Get company by user's company ID
	company, err := h.databaseProvider.GetCompany(c.Request.Context(), user.CompanyID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Company not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: gin.H{
			"company": gin.H{
				"id":          company.ID,
				"name":        company.Name,
				"domain":      company.Domain,
				"color_theme": company.ColorTheme,
				"logo_url":    company.LogoURL,
				"status":      company.Status,
				"trial_ends_at": company.TrialEndsAt,
				"created_at":    company.CreatedAt,
				"updated_at":    company.UpdatedAt,
				"onboarded":     company.Onboarded,
			},
		},
	})
}

// UpdateCompany handles company updates
func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	var req models.UpdateCompanyRequest
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

	// Get company
	company, err := h.databaseProvider.GetCompany(c.Request.Context(), user.CompanyID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Company not found",
		})
		return
	}

	// Check if user is admin
	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "Only admins can update company settings",
		})
		return
	}

	// Update company fields
	if req.Name != "" {
		company.Name = req.Name
	}
	if req.Domain != "" {
		// Check if new domain is already taken
		if req.Domain != company.Domain {
			existingCompany, err := h.databaseProvider.GetCompanyByDomain(c.Request.Context(), req.Domain)
			if err == nil && existingCompany != nil {
				c.JSON(http.StatusConflict, models.APIResponse{
					Success: false,
					Error:   "Domain already taken",
				})
				return
			}
			company.Domain = req.Domain
		}
	}
	if req.ColorTheme != "" {
		company.ColorTheme = req.ColorTheme
	}
	if req.LogoURL != "" {
		company.LogoURL = req.LogoURL
	}

	// Save updated company
	if err := h.databaseProvider.UpdateCompany(c.Request.Context(), company); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to update company",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Company updated successfully",
		Data: gin.H{
			"company": gin.H{
				"id":          company.ID,
				"name":        company.Name,
				"domain":      company.Domain,
				"color_theme": company.ColorTheme,
				"logo_url":    company.LogoURL,
				"status":      company.Status,
				"updated_at":  company.UpdatedAt,
			},
		},
	})
}

// DeleteCompany handles company deletion
func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
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
			Error:   "Only admins can delete company",
		})
		return
	}

	// Delete company
	if err := h.databaseProvider.DeleteCompany(c.Request.Context(), user.CompanyID); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to delete company",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Company deleted successfully",
	})
}

// ListCompanies handles listing companies (admin only)
func (h *CompanyHandler) ListCompanies(c *gin.Context) {
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

	// Check if user is admin (assuming admin users can see all companies)
	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "Access denied",
		})
		return
	}

	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Get companies
	companies, err := h.databaseProvider.ListCompanies(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to get companies",
		})
		return
	}

	// Convert to response format
	var companyList []gin.H
	for _, company := range companies {
		companyList = append(companyList, gin.H{
			"id":          company.ID,
			"name":        company.Name,
			"domain":      company.Domain,
			"color_theme": company.ColorTheme,
			"status":      company.Status,
			"admin_user_id": company.AdminUserID,
			"created_at":    company.CreatedAt,
			"updated_at":    company.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: gin.H{
			"companies": companyList,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
			},
		},
	})
}

// GetCompanyStats handles getting company statistics
func (h *CompanyHandler) GetCompanyStats(c *gin.Context) {
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

	// Get user count
	userCount, err := h.databaseProvider.CountUsersByCompany(c.Request.Context(), user.CompanyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to get user count",
		})
		return
	}

	// Get company
	company, err := h.databaseProvider.GetCompany(c.Request.Context(), user.CompanyID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Company not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: gin.H{
			"stats": gin.H{
				"total_users": userCount,
				"company_status": company.Status,
				"trial_ends_at": company.TrialEndsAt,
				"onboarded": company.Onboarded,
			},
		},
	})
}

