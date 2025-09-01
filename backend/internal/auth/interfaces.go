package auth

import (
	"context"
	"time"
)

// User represents a user in the system
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Picture      string    `json:"picture,omitempty"`
	CompanyID    string    `json:"company_id"`
	Role         UserRole  `json:"role"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	LastLoginAt  time.Time `json:"last_login_at,omitempty"`
	OnboardedAt  time.Time `json:"onboarded_at,omitempty"`
	Onboarded    bool      `json:"onboarded"`
}

// UserRole represents the role of a user
type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleUser   UserRole = "user"
	RoleGuest  UserRole = "guest"
)

// AuthProvider defines the interface for authentication providers
type AuthProvider interface {
	// Authenticate authenticates a user with the given credentials
	Authenticate(ctx context.Context, email, password string) (*User, error)
	
	// AuthenticateWithToken authenticates a user with a token (JWT, OAuth, etc.)
	AuthenticateWithToken(ctx context.Context, token string) (*User, error)
	
	// Register registers a new user
	Register(ctx context.Context, email, password, name string) (*User, error)
	
	// RegisterWithSocial registers a new user with social login
	RegisterWithSocial(ctx context.Context, provider, token string) (*User, error)
	
	// GetUser retrieves a user by ID
	GetUser(ctx context.Context, userID string) (*User, error)
	
	// GetUserByEmail retrieves a user by email
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	
	// UpdateUser updates user information
	UpdateUser(ctx context.Context, user *User) error
	
	// DeleteUser deletes a user
	DeleteUser(ctx context.Context, userID string) error
	
	// GenerateToken generates a JWT token for a user
	GenerateToken(user *User) (string, error)
	
	// ValidateToken validates a JWT token and returns the user
	ValidateToken(token string) (*User, error)
	
	// RefreshToken refreshes a JWT token
	RefreshToken(refreshToken string) (string, error)
	
	// SendInvitation sends an invitation email to a user
	SendInvitation(ctx context.Context, email, companyID string, invitedBy string) error
	
	// ActivateUser activates a user account
	ActivateUser(ctx context.Context, activationToken string) error
	
	// ResetPassword sends a password reset email
	ResetPassword(ctx context.Context, email string) error
	
	// ChangePassword changes a user's password
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
}

// AuthConfig holds configuration for authentication providers
type AuthConfig struct {
	Provider     string `json:"provider"`      // "auth0", "google", "custom"
	Domain       string `json:"domain"`        // Auth0 domain or OAuth provider domain
	ClientID     string `json:"client_id"`     // OAuth client ID
	ClientSecret string `json:"client_secret"` // OAuth client secret
	JWTSecret    string `json:"jwt_secret"`    // JWT signing secret
	RedirectURL  string `json:"redirect_url"`  // OAuth redirect URL
}

// AuthFactory creates authentication providers
type AuthFactory interface {
	CreateProvider(config AuthConfig) (AuthProvider, error)
}

// DefaultAuthFactory implements AuthFactory
type DefaultAuthFactory struct{}

// CreateProvider creates an auth provider based on configuration
func (f *DefaultAuthFactory) CreateProvider(config AuthConfig) (AuthProvider, error) {
	switch config.Provider {
	case "auth0":
		return NewAuth0Provider(config)
	case "google":
		return NewGoogleProvider(config)
	case "custom":
		return NewCustomProvider(config)
	default:
		return NewCustomProvider(config)
	}
}

