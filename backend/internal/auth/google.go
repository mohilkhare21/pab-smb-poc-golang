package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GoogleProvider implements AuthProvider for Google OAuth
type GoogleProvider struct {
	config     AuthConfig
	httpClient *http.Client
}

// NewGoogleProvider creates a new Google OAuth provider
func NewGoogleProvider(config AuthConfig) (*GoogleProvider, error) {
	return &GoogleProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

// GoogleUserInfo represents user information from Google
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// Authenticate authenticates a user with email and password
func (g *GoogleProvider) Authenticate(ctx context.Context, email, password string) (*User, error) {
	// Google OAuth doesn't support direct username/password authentication
	return nil, fmt.Errorf("direct authentication not supported with Google OAuth, use AuthenticateWithToken")
}

// AuthenticateWithToken authenticates a user with a Google OAuth token
func (g *GoogleProvider) AuthenticateWithToken(ctx context.Context, token string) (*User, error) {
	// Get user info from Google
	userInfo, err := g.getGoogleUserInfo(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get Google user info: %w", err)
	}

	// Create user from Google info
	user := &User{
		ID:        userInfo.ID,
		Email:     userInfo.Email,
		Name:      userInfo.Name,
		Picture:   userInfo.Picture,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return user, nil
}

// Register registers a new user with Google OAuth
func (g *GoogleProvider) Register(ctx context.Context, email, password, name string) (*User, error) {
	// For Google OAuth, registration happens through the OAuth flow
	// This method is not typically used with Google OAuth
	return nil, fmt.Errorf("registration not supported with Google OAuth, use RegisterWithSocial")
}

// RegisterWithSocial registers a new user with Google OAuth
func (g *GoogleProvider) RegisterWithSocial(ctx context.Context, provider, token string) (*User, error) {
	if provider != "google" {
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	// Get user info from Google
	userInfo, err := g.getGoogleUserInfo(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get Google user info: %w", err)
	}

	// Create user from Google info
	user := &User{
		ID:        userInfo.ID,
		Email:     userInfo.Email,
		Name:      userInfo.Name,
		Picture:   userInfo.Picture,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (g *GoogleProvider) GetUser(ctx context.Context, userID string) (*User, error) {
	// For Google OAuth, we'd need to store user info in our database
	// For POC, we'll return a mock user
	user := &User{
		ID:       userID,
		Email:    "user@example.com",
		Name:     "Google User",
		IsActive: true,
	}
	return user, nil
}

// GetUserByEmail retrieves a user by email
func (g *GoogleProvider) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	// For Google OAuth, we'd need to store user info in our database
	// For POC, we'll return a mock user
	user := &User{
		ID:       uuid.New().String(),
		Email:    email,
		Name:     "Google User",
		IsActive: true,
	}
	return user, nil
}

// UpdateUser updates user information
func (g *GoogleProvider) UpdateUser(ctx context.Context, user *User) error {
	// For Google OAuth, user profile updates are handled by Google
	// We can only update our local user data
	return nil
}

// DeleteUser deletes a user
func (g *GoogleProvider) DeleteUser(ctx context.Context, userID string) error {
	// For Google OAuth, user deletion is handled by Google
	// We can only delete our local user data
	return nil
}

// GenerateToken generates a JWT token for a user
func (g *GoogleProvider) GenerateToken(user *User) (string, error) {
	// Create JWT claims
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"name":  user.Name,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(24 * time.Hour).Unix(), // 24 hour expiry
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Sign token
	tokenString, err := token.SignedString([]byte(g.config.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the user
func (g *GoogleProvider) ValidateToken(tokenString string) (*User, error) {
	// Parse and validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(g.config.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Create user from claims
	user := &User{
		ID:       claims["sub"].(string),
		Email:    claims["email"].(string),
		Name:     claims["name"].(string),
		IsActive: true,
	}

	return user, nil
}

// RefreshToken refreshes a JWT token
func (g *GoogleProvider) RefreshToken(refreshToken string) (string, error) {
	// For simplicity, we'll just generate a new token
	// In a real implementation, you'd validate the refresh token and generate a new access token
	return "", fmt.Errorf("refresh token not implemented")
}

// SendInvitation sends an invitation email to a user
func (g *GoogleProvider) SendInvitation(ctx context.Context, email, companyID string, invitedBy string) error {
	// For Google OAuth, invitations would typically be sent via email
	// For POC, we'll just log the invitation
	fmt.Printf("Invitation sent to %s for company %s by %s\n", email, companyID, invitedBy)
	return nil
}

// ActivateUser activates a user account
func (g *GoogleProvider) ActivateUser(ctx context.Context, activationToken string) error {
	// For Google OAuth, user activation is handled by Google
	return nil
}

// ResetPassword sends a password reset email
func (g *GoogleProvider) ResetPassword(ctx context.Context, email string) error {
	// For Google OAuth, password reset is handled by Google
	fmt.Printf("Password reset requested for %s (handled by Google)\n", email)
	return nil
}

// ChangePassword changes a user's password
func (g *GoogleProvider) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	// For Google OAuth, password changes are handled by Google
	return fmt.Errorf("password changes are handled by Google")
}

// getGoogleUserInfo gets user information from Google
func (g *GoogleProvider) getGoogleUserInfo(ctx context.Context, accessToken string) (*GoogleUserInfo, error) {
	// Make request to Google's userinfo endpoint
	req, err := http.NewRequestWithContext(ctx, "GET", 
		"https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: %d", resp.StatusCode)
	}

	// Parse response
	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}
