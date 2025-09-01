package models

import (
	"time"
)

// Company represents a company/tenant in the system
type Company struct {
	ID              string    `json:"id" firestore:"id"`
	Name            string    `json:"name" firestore:"name"`
	Domain          string    `json:"domain" firestore:"domain"`
	ColorTheme      string    `json:"color_theme" firestore:"color_theme"`
	LogoURL         string    `json:"logo_url,omitempty" firestore:"logo_url,omitempty"`
	AdminUserID     string    `json:"admin_user_id" firestore:"admin_user_id"`
	SubscriptionID  string    `json:"subscription_id,omitempty" firestore:"subscription_id,omitempty"`
	Status          string    `json:"status" firestore:"status"` // "active", "trial", "suspended", "cancelled"
	TrialEndsAt     time.Time `json:"trial_ends_at,omitempty" firestore:"trial_ends_at,omitempty"`
	CreatedAt       time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" firestore:"updated_at"`
	OnboardedAt     time.Time `json:"onboarded_at,omitempty" firestore:"onboarded_at,omitempty"`
	Onboarded       bool      `json:"onboarded" firestore:"onboarded"`
	// New fields for configuration status tracking
	SetupCompleted  bool      `json:"setup_completed" firestore:"setup_completed"`
	SetupCompletedAt time.Time `json:"setup_completed_at,omitempty" firestore:"setup_completed_at,omitempty"`
	// Configuration status for each feature
	WebsiteSecurityConfigured bool `json:"website_security_configured" firestore:"website_security_configured"`
	MalwareSecurityConfigured bool `json:"malware_security_configured" firestore:"malware_security_configured"`
	DataControlsConfigured    bool `json:"data_controls_configured" firestore:"data_controls_configured"`
	ReportingConfigured       bool `json:"reporting_configured" firestore:"reporting_configured"`
	BrowserCustomized         bool `json:"browser_customized" firestore:"browser_customized"`
	SubscriptionActive        bool `json:"subscription_active" firestore:"subscription_active"`
	UsersInvited              bool `json:"users_invited" firestore:"users_invited"`
	DownloadReady             bool `json:"download_ready" firestore:"download_ready"`
}

// User represents a user in the system
type User struct {
	ID           string    `json:"id" firestore:"id"`
	Email        string    `json:"email" firestore:"email"`
	Name         string    `json:"name" firestore:"name"`
	Picture      string    `json:"picture,omitempty" firestore:"picture,omitempty"`
	CompanyID    string    `json:"company_id" firestore:"company_id"`
	Role         UserRole  `json:"role" firestore:"role"`
	IsActive     bool      `json:"is_active" firestore:"is_active"`
	CreatedAt    time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" firestore:"updated_at"`
	LastLoginAt  time.Time `json:"last_login_at,omitempty" firestore:"last_login_at,omitempty"`
	OnboardedAt  time.Time `json:"onboarded_at,omitempty" firestore:"onboarded_at,omitempty"`
	Onboarded    bool      `json:"onboarded" firestore:"onboarded"`
	// New fields for invitation status
	InvitationStatus string    `json:"invitation_status" firestore:"invitation_status"` // "invited", "active", "pending"
	InvitedAt        time.Time `json:"invited_at,omitempty" firestore:"invited_at,omitempty"`
	ActivatedAt      time.Time `json:"activated_at,omitempty" firestore:"activated_at,omitempty"`
}

// UserRole represents the role of a user
type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleUser   UserRole = "user"
	RoleGuest  UserRole = "guest"
)

// Invitation represents a user invitation
type Invitation struct {
	ID           string    `json:"id" firestore:"id"`
	Email        string    `json:"email" firestore:"email"`
	CompanyID    string    `json:"company_id" firestore:"company_id"`
	InvitedBy    string    `json:"invited_by" firestore:"invited_by"`
	Token        string    `json:"token" firestore:"token"`
	Status       string    `json:"status" firestore:"status"` // "pending", "accepted", "expired", "sent"
	ExpiresAt    time.Time `json:"expires_at" firestore:"expires_at"`
	CreatedAt    time.Time `json:"created_at" firestore:"created_at"`
	AcceptedAt   time.Time `json:"accepted_at,omitempty" firestore:"accepted_at,omitempty"`
	SentAt       time.Time `json:"sent_at,omitempty" firestore:"sent_at,omitempty"`
	// New fields for tracking
	SentCount    int       `json:"sent_count" firestore:"sent_count"`
	LastSentAt   time.Time `json:"last_sent_at,omitempty" firestore:"last_sent_at,omitempty"`
}

// BrowserShortcut represents a browser shortcut for a company
type BrowserShortcut struct {
	ID          string `json:"id" firestore:"id"`
	CompanyID   string `json:"company_id" firestore:"company_id"`
	Name        string `json:"name" firestore:"name"`
	URL         string `json:"url" firestore:"url"`
	Icon        string `json:"icon,omitempty" firestore:"icon,omitempty"`
	Description string `json:"description,omitempty" firestore:"description,omitempty"`
	Order       int    `json:"order" firestore:"order"`
	IsActive    bool   `json:"is_active" firestore:"is_active"`
	// New fields for shortcut management
	IsSuggested bool   `json:"is_suggested" firestore:"is_suggested"` // Whether this was auto-generated
	Category    string `json:"category" firestore:"category"`        // "company", "suggested", "custom"
	Source      string `json:"source,omitempty" firestore:"source,omitempty"` // How this shortcut was added
}

// Subscription represents a subscription for a company
type Subscription struct {
	ID                 string    `json:"id" firestore:"id"`
	CompanyID          string    `json:"company_id" firestore:"company_id"`
	StripeID           string    `json:"stripe_id" firestore:"stripe_id"`
	Plan               string    `json:"plan" firestore:"plan"`
	Status             string    `json:"status" firestore:"status"`
	CurrentPeriodStart time.Time `json:"current_period_start" firestore:"current_period_start"`
	CurrentPeriodEnd   time.Time `json:"current_period_end" firestore:"current_period_end"`
	TrialStart         time.Time `json:"trial_start,omitempty" firestore:"trial_start,omitempty"`
	TrialEnd           time.Time `json:"trial_end,omitempty" firestore:"trial_end,omitempty"`
	CreatedAt          time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" firestore:"updated_at"`
	// New fields for user management
	MaxUsers           int       `json:"max_users" firestore:"max_users"`
	ActiveUsers        int       `json:"active_users" firestore:"active_users"`
	InvitedUsers       int       `json:"invited_users" firestore:"invited_users"`
	IsTrialActive      bool      `json:"is_trial_active" firestore:"is_trial_active"`
	TrialDaysRemaining int       `json:"trial_days_remaining" firestore:"trial_days_remaining"`
}

// CompanySetupProgress represents the setup progress for a company
type CompanySetupProgress struct {
	CompanyID                    string    `json:"company_id" firestore:"company_id"`
	Step                        string    `json:"step" firestore:"step"` // "domain", "customization", "invitations", "subscription", "complete"
	Progress                    int       `json:"progress" firestore:"progress"` // 0-100
	DomainProvided              bool      `json:"domain_provided" firestore:"domain_provided"`
	CustomizationCompleted      bool      `json:"customization_completed" firestore:"customization_completed"`
	InvitationsSent             bool      `json:"invitations_sent" firestore:"invitations_sent"`
	SubscriptionStarted         bool      `json:"subscription_started" firestore:"subscription_started"`
	SetupCompleted              bool      `json:"setup_completed" firestore:"setup_completed"`
	LastUpdated                 time.Time `json:"last_updated" firestore:"last_updated"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

// CompanyCreateRequest represents a company creation request
type CompanyCreateRequest struct {
	Name       string `json:"name" binding:"required"`
	Domain     string `json:"domain" binding:"required"`
	ColorTheme string `json:"color_theme" binding:"required"`
}

// InviteUserRequest represents a user invitation request
type InviteUserRequest struct {
	Emails []string `json:"emails" binding:"required"`
}

// UpdateCompanyRequest represents a company update request
type UpdateCompanyRequest struct {
	Name       string `json:"name,omitempty"`
	Domain     string `json:"domain,omitempty"`
	ColorTheme string `json:"color_theme,omitempty"`
	LogoURL    string `json:"logo_url,omitempty"`
}

// UpdateUserRequest represents a user update request
type UpdateUserRequest struct {
	Name     string   `json:"name,omitempty"`
	Role     UserRole `json:"role,omitempty"`
	IsActive *bool    `json:"is_active,omitempty"`
}

// BrowserShortcutRequest represents a browser shortcut request
type BrowserShortcutRequest struct {
	Name        string `json:"name" binding:"required"`
	URL         string `json:"url" binding:"required,url"`
	Icon        string `json:"icon,omitempty"`
	Description string `json:"description,omitempty"`
	Order       int    `json:"order"`
	Category    string `json:"category,omitempty"`
}

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page  int `json:"page" form:"page"`
	Limit int `json:"limit" form:"limit"`
}

// PaginationResponse represents paginated response
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

// UserContext represents user context in requests
type UserContext struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	CompanyID string `json:"company_id"`
	Role      string `json:"role"`
}

// New request/response models for enhanced functionality

// CompanySetupRequest represents company setup request
type CompanySetupRequest struct {
	Domain     string `json:"domain" binding:"required"`
	ColorTheme string `json:"color_theme" binding:"required"`
	LogoURL    string `json:"logo_url,omitempty"`
}

// CompanyStatsResponse represents company statistics
type CompanyStatsResponse struct {
	TotalUsers       int `json:"total_users"`
	ActiveUsers      int `json:"active_users"`
	InvitedUsers     int `json:"invited_users"`
	PendingInvitations int `json:"pending_invitations"`
	MaxUsers         int `json:"max_users"`
	SetupProgress    int `json:"setup_progress"`
	ConfigurationStatus map[string]bool `json:"configuration_status"`
}

// GenerateShortcutsRequest represents request to generate shortcuts
type GenerateShortcutsRequest struct {
	Domain string `json:"domain" binding:"required"`
}

// NudgeUsersRequest represents request to nudge users
type NudgeUsersRequest struct {
	UserIDs []string `json:"user_ids" binding:"required"`
	Message string   `json:"message,omitempty"`
}

// ConfigurationStatusRequest represents configuration status update
type ConfigurationStatusRequest struct {
	Feature string `json:"feature" binding:"required"`
	Status  bool   `json:"status"`
}

// DownloadInfoResponse represents download information
type DownloadInfoResponse struct {
	DownloadURL    string `json:"download_url"`
	Version        string `json:"version"`
	ReleaseDate    string `json:"release_date"`
	FileSize       string `json:"file_size"`
	SupportedOS    []string `json:"supported_os"`
	InstallationInstructions string `json:"installation_instructions"`
}
