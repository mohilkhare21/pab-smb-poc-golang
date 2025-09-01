package database

import (
	"context"
	"time"
)

// Company represents a company/tenant in the system
type Company struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Domain          string    `json:"domain"`
	ColorTheme      string    `json:"color_theme"`
	LogoURL         string    `json:"logo_url,omitempty"`
	AdminUserID     string    `json:"admin_user_id"`
	SubscriptionID  string    `json:"subscription_id,omitempty"`
	Status          string    `json:"status"` // "active", "trial", "suspended", "cancelled"
	TrialEndsAt     time.Time `json:"trial_ends_at,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	OnboardedAt     time.Time `json:"onboarded_at,omitempty"`
	Onboarded       bool      `json:"onboarded"`
	// New fields for configuration status tracking
	SetupCompleted  bool      `json:"setup_completed"`
	SetupCompletedAt time.Time `json:"setup_completed_at,omitempty"`
	// Configuration status for each feature
	WebsiteSecurityConfigured bool `json:"website_security_configured"`
	MalwareSecurityConfigured bool `json:"malware_security_configured"`
	DataControlsConfigured    bool `json:"data_controls_configured"`
	ReportingConfigured       bool `json:"reporting_configured"`
	BrowserCustomized         bool `json:"browser_customized"`
	SubscriptionActive        bool `json:"subscription_active"`
	UsersInvited              bool `json:"users_invited"`
	DownloadReady             bool `json:"download_ready"`
}

// User represents a user in the database
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Picture      string    `json:"picture,omitempty"`
	CompanyID    string    `json:"company_id"`
	Role         string    `json:"role"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	LastLoginAt  time.Time `json:"last_login_at,omitempty"`
	OnboardedAt  time.Time `json:"onboarded_at,omitempty"`
	Onboarded    bool      `json:"onboarded"`
	// New fields for invitation status
	InvitationStatus string    `json:"invitation_status"` // "invited", "active", "pending"
	InvitedAt        time.Time `json:"invited_at,omitempty"`
	ActivatedAt      time.Time `json:"activated_at,omitempty"`
}

// Invitation represents a user invitation
type Invitation struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	CompanyID    string    `json:"company_id"`
	InvitedBy    string    `json:"invited_by"`
	Token        string    `json:"token"`
	Status       string    `json:"status"` // "pending", "accepted", "expired", "sent"
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	AcceptedAt   time.Time `json:"accepted_at,omitempty"`
	SentAt       time.Time `json:"sent_at,omitempty"`
	// New fields for tracking
	SentCount    int       `json:"sent_count"`
	LastSentAt   time.Time `json:"last_sent_at,omitempty"`
}

// BrowserShortcut represents a browser shortcut for a company
type BrowserShortcut struct {
	ID          string `json:"id"`
	CompanyID   string `json:"company_id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Icon        string `json:"icon,omitempty"`
	Description string `json:"description,omitempty"`
	Order       int    `json:"order"`
	IsActive    bool   `json:"is_active"`
	// New fields for shortcut management
	IsSuggested bool   `json:"is_suggested"` // Whether this was auto-generated
	Category    string `json:"category"`    // "company", "suggested", "custom"
	Source      string `json:"source,omitempty"` // How this shortcut was added
}

// Subscription represents a subscription for a company
type Subscription struct {
	ID                 string    `json:"id"`
	CompanyID          string    `json:"company_id"`
	StripeID           string    `json:"stripe_id"`
	Plan               string    `json:"plan"`
	Status             string    `json:"status"`
	CurrentPeriodStart time.Time `json:"current_period_start"`
	CurrentPeriodEnd   time.Time `json:"current_period_end"`
	TrialStart         time.Time `json:"trial_start,omitempty"`
	TrialEnd           time.Time `json:"trial_end,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	// New fields for user management
	MaxUsers           int       `json:"max_users"`
	ActiveUsers        int       `json:"active_users"`
	InvitedUsers       int       `json:"invited_users"`
	IsTrialActive      bool      `json:"is_trial_active"`
	TrialDaysRemaining int       `json:"trial_days_remaining"`
}

// CompanySetupProgress represents the setup progress for a company
type CompanySetupProgress struct {
	CompanyID                    string    `json:"company_id"`
	Step                        string    `json:"step"` // "domain", "customization", "invitations", "subscription", "complete"
	Progress                    int       `json:"progress"` // 0-100
	DomainProvided              bool      `json:"domain_provided"`
	CustomizationCompleted      bool      `json:"customization_completed"`
	InvitationsSent             bool      `json:"invitations_sent"`
	SubscriptionStarted         bool      `json:"subscription_started"`
	SetupCompleted              bool      `json:"setup_completed"`
	LastUpdated                 time.Time `json:"last_updated"`
}

// DatabaseProvider defines the interface for database providers
type DatabaseProvider interface {
	// Company operations
	CreateCompany(ctx context.Context, company *Company) error
	GetCompany(ctx context.Context, companyID string) (*Company, error)
	GetCompanyByDomain(ctx context.Context, domain string) (*Company, error)
	UpdateCompany(ctx context.Context, company *Company) error
	DeleteCompany(ctx context.Context, companyID string) error
	ListCompanies(ctx context.Context, limit, offset int) ([]*Company, error)
	
	// User operations
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, userID string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUsersByCompany(ctx context.Context, companyID string) ([]*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, userID string) error
	CountUsersByCompany(ctx context.Context, companyID string) (int, error)
	
	// Enhanced user operations for invitation status
	GetInvitedUsersByCompany(ctx context.Context, companyID string) ([]*User, error)
	GetActiveUsersByCompany(ctx context.Context, companyID string) ([]*User, error)
	CountInvitedUsersByCompany(ctx context.Context, companyID string) (int, error)
	CountActiveUsersByCompany(ctx context.Context, companyID string) (int, error)
	UpdateUserInvitationStatus(ctx context.Context, userID string, status string) error
	
	// Invitation operations
	CreateInvitation(ctx context.Context, invitation *Invitation) error
	GetInvitation(ctx context.Context, invitationID string) (*Invitation, error)
	GetInvitationByToken(ctx context.Context, token string) (*Invitation, error)
	GetInvitationsByCompany(ctx context.Context, companyID string) ([]*Invitation, error)
	UpdateInvitation(ctx context.Context, invitation *Invitation) error
	DeleteInvitation(ctx context.Context, invitationID string) error
	DeleteExpiredInvitations(ctx context.Context) error
	
	// Enhanced invitation operations
	GetPendingInvitationsByCompany(ctx context.Context, companyID string) ([]*Invitation, error)
	CountPendingInvitationsByCompany(ctx context.Context, companyID string) (int, error)
	UpdateInvitationSentStatus(ctx context.Context, invitationID string, sentAt time.Time) error
	ResendInvitation(ctx context.Context, invitationID string) error
	
	// Browser shortcut operations
	CreateBrowserShortcut(ctx context.Context, shortcut *BrowserShortcut) error
	GetBrowserShortcut(ctx context.Context, shortcutID string) (*BrowserShortcut, error)
	GetBrowserShortcutsByCompany(ctx context.Context, companyID string) ([]*BrowserShortcut, error)
	UpdateBrowserShortcut(ctx context.Context, shortcut *BrowserShortcut) error
	DeleteBrowserShortcut(ctx context.Context, shortcutID string) error
	DeleteBrowserShortcutsByCompany(ctx context.Context, companyID string) error
	
	// Enhanced shortcut operations
	GetSuggestedShortcutsByCompany(ctx context.Context, companyID string) ([]*BrowserShortcut, error)
	GetCustomShortcutsByCompany(ctx context.Context, companyID string) ([]*BrowserShortcut, error)
	GenerateShortcutsForDomain(ctx context.Context, companyID string, domain string) error
	
	// Subscription operations
	CreateSubscription(ctx context.Context, subscription *Subscription) error
	GetSubscription(ctx context.Context, subscriptionID string) (*Subscription, error)
	GetSubscriptionByCompany(ctx context.Context, companyID string) (*Subscription, error)
	UpdateSubscription(ctx context.Context, subscription *Subscription) error
	DeleteSubscription(ctx context.Context, subscriptionID string) error
	
	// Enhanced subscription operations
	UpdateSubscriptionUserCounts(ctx context.Context, subscriptionID string, activeUsers, invitedUsers int) error
	GetSubscriptionStats(ctx context.Context, companyID string) (*Subscription, error)
	
	// Setup progress operations
	CreateSetupProgress(ctx context.Context, progress *CompanySetupProgress) error
	GetSetupProgress(ctx context.Context, companyID string) (*CompanySetupProgress, error)
	UpdateSetupProgress(ctx context.Context, progress *CompanySetupProgress) error
	UpdateSetupStep(ctx context.Context, companyID string, step string, progress int) error
	
	// Configuration status operations
	UpdateCompanyConfigurationStatus(ctx context.Context, companyID string, feature string, status bool) error
	GetCompanyConfigurationStatus(ctx context.Context, companyID string) (map[string]bool, error)
	
	// Transaction operations
	BeginTransaction(ctx context.Context) (Transaction, error)
	
	// Health check
	Ping(ctx context.Context) error
	Close() error
}

// Transaction represents a database transaction
type Transaction interface {
	Commit() error
	Rollback() error
}

// DatabaseConfig holds configuration for database providers
type DatabaseConfig struct {
	Provider    string `json:"provider"`     // "firestore", "postgres", "mysql"
	ProjectID   string `json:"project_id"`   // Firebase project ID or database name
	Host        string `json:"host"`         // Database host
	Port        int    `json:"port"`         // Database port
	Username    string `json:"username"`     // Database username
	Password    string `json:"password"`     // Database password
	Database    string `json:"database"`     // Database name
	SSLMode     string `json:"ssl_mode"`     // SSL mode for SQL databases
	Credentials string `json:"credentials"`  // Path to service account key for Firestore
}

// DatabaseFactory creates database providers
type DatabaseFactory interface {
	CreateProvider(config DatabaseConfig) (DatabaseProvider, error)
}

// DefaultDatabaseFactory implements DatabaseFactory
type DefaultDatabaseFactory struct{}

// CreateProvider creates a database provider based on configuration
func (f *DefaultDatabaseFactory) CreateProvider(config DatabaseConfig) (DatabaseProvider, error) {
	switch config.Provider {
	case "firestore":
		return NewFirestoreProvider(config)
	case "postgres":
		return NewPostgresProvider(config)
	case "mysql":
		return NewMySQLProvider(config)
	default:
		return NewFirestoreProvider(config)
	}
}
