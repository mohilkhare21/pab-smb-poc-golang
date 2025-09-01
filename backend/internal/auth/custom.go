package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// CustomProvider implements AuthProvider for custom authentication
type CustomProvider struct {
	config AuthConfig
	users  map[string]*User // In-memory user store for POC
}

// NewCustomProvider creates a new custom auth provider
func NewCustomProvider(config AuthConfig) (*CustomProvider, error) {
	return &CustomProvider{
		config: config,
		users:  make(map[string]*User),
	}, nil
}

// Authenticate authenticates a user with email and password
func (c *CustomProvider) Authenticate(ctx context.Context, email, password string) (*User, error) {
	// In a real implementation, you'd hash the password and compare with stored hash
	// For POC, we'll use a simple approach
	for _, user := range c.users {
		if user.Email == email {
			// In real implementation, verify password hash
			return user, nil
		}
	}
	return nil, fmt.Errorf("invalid credentials")
}

// AuthenticateWithToken authenticates a user with a token
func (c *CustomProvider) AuthenticateWithToken(ctx context.Context, token string) (*User, error) {
	// Validate the JWT token
	claims, err := c.validateJWT(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Extract user information from claims
	userID := claims["sub"].(string)
	
	// Get user from memory store
	user, exists := c.users[userID]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

// Register registers a new user
func (c *CustomProvider) Register(ctx context.Context, email, password, name string) (*User, error) {
	// Check if user already exists
	for _, user := range c.users {
		if user.Email == email {
			return nil, fmt.Errorf("user already exists")
		}
	}

	// Create new user
	user := &User{
		ID:        uuid.New().String(),
		Email:     email,
		Name:      name,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Store user in memory
	c.users[user.ID] = user

	return user, nil
}

// RegisterWithSocial registers a new user with social login
func (c *CustomProvider) RegisterWithSocial(ctx context.Context, provider, token string) (*User, error) {
	// For POC, we'll create a user with basic info
	user := &User{
		ID:        uuid.New().String(),
		Email:     fmt.Sprintf("user-%s@example.com", uuid.New().String()[:8]),
		Name:      fmt.Sprintf("User from %s", provider),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Store user in memory
	c.users[user.ID] = user

	return user, nil
}

// GetUser retrieves a user by ID
func (c *CustomProvider) GetUser(ctx context.Context, userID string) (*User, error) {
	user, exists := c.users[userID]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// GetUserByEmail retrieves a user by email
func (c *CustomProvider) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	for _, user := range c.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

// UpdateUser updates user information
func (c *CustomProvider) UpdateUser(ctx context.Context, user *User) error {
	_, exists := c.users[user.ID]
	if !exists {
		return fmt.Errorf("user not found")
	}

	user.UpdatedAt = time.Now()
	c.users[user.ID] = user
	return nil
}

// DeleteUser deletes a user
func (c *CustomProvider) DeleteUser(ctx context.Context, userID string) error {
	_, exists := c.users[userID]
	if !exists {
		return fmt.Errorf("user not found")
	}

	delete(c.users, userID)
	return nil
}

// GenerateToken generates a JWT token for a user
func (c *CustomProvider) GenerateToken(user *User) (string, error) {
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
	tokenString, err := token.SignedString([]byte(c.config.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the user
func (c *CustomProvider) ValidateToken(tokenString string) (*User, error) {
	// Parse and validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(c.config.JWTSecret), nil
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
func (c *CustomProvider) RefreshToken(refreshToken string) (string, error) {
	// For simplicity, we'll just generate a new token
	// In a real implementation, you'd validate the refresh token and generate a new access token
	return "", fmt.Errorf("refresh token not implemented")
}

// SendInvitation sends an invitation email to a user
func (c *CustomProvider) SendInvitation(ctx context.Context, email, companyID string, invitedBy string) error {
	// For POC, we'll just log the invitation
	// In a real implementation, you'd send an email
	fmt.Printf("Invitation sent to %s for company %s by %s\n", email, companyID, invitedBy)
	return nil
}

// ActivateUser activates a user account
func (c *CustomProvider) ActivateUser(ctx context.Context, activationToken string) error {
	// For POC, we'll just return success
	// In a real implementation, you'd validate the activation token and activate the user
	return nil
}

// ResetPassword sends a password reset email
func (c *CustomProvider) ResetPassword(ctx context.Context, email string) error {
	// For POC, we'll just log the password reset request
	// In a real implementation, you'd send a password reset email
	fmt.Printf("Password reset requested for %s\n", email)
	return nil
}

// ChangePassword changes a user's password
func (c *CustomProvider) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	// For POC, we'll just return success
	// In a real implementation, you'd validate the old password and update with the new password hash
	return nil
}

// validateJWT validates a JWT token
func (c *CustomProvider) validateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(c.config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// generateRandomToken generates a random token for invitations
func (c *CustomProvider) generateRandomToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
