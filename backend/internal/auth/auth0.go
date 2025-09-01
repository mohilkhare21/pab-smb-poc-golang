package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Auth0Provider implements AuthProvider for Auth0
type Auth0Provider struct {
	config     AuthConfig
	httpClient *http.Client
}

// NewAuth0Provider creates a new Auth0 provider
func NewAuth0Provider(config AuthConfig) (*Auth0Provider, error) {
	return &Auth0Provider{
		config: config,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

// Auth0User represents a user from Auth0
type Auth0User struct {
	UserID      string `json:"user_id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Picture     string `json:"picture"`
	EmailVerified bool `json:"email_verified"`
	UpdatedAt   string `json:"updated_at"`
}

// Auth0Token represents an Auth0 token response
type Auth0Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope"`
}

// Auth0ManagementToken represents an Auth0 management API token
type Auth0ManagementToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// Authenticate authenticates a user with email and password
func (a *Auth0Provider) Authenticate(ctx context.Context, email, password string) (*User, error) {
	// Auth0 doesn't support direct username/password authentication via API
	// This would typically be handled by the frontend using Auth0's Universal Login
	return nil, fmt.Errorf("direct authentication not supported with Auth0, use AuthenticateWithToken")
}

// AuthenticateWithToken authenticates a user with an Auth0 token
func (a *Auth0Provider) AuthenticateWithToken(ctx context.Context, token string) (*User, error) {
	// Validate the JWT token
	claims, err := a.validateJWT(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Extract user information from claims
	user := &User{
		ID:        claims["sub"].(string),
		Email:     claims["email"].(string),
		Name:      claims["name"].(string),
		Picture:   claims["picture"].(string),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Check if user exists in our database
	// This would typically be done by calling the database provider
	// For now, we'll return the user from the token

	return user, nil
}

// Register registers a new user with Auth0
func (a *Auth0Provider) Register(ctx context.Context, email, password, name string) (*User, error) {
	// Get management API token
	mgmtToken, err := a.getManagementToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get management token: %w", err)
	}

	// Create user in Auth0
	userID := uuid.New().String()
	auth0User := map[string]interface{}{
		"user_id":      userID,
		"email":        email,
		"password":     password,
		"name":         name,
		"connection":   "Username-Password-Authentication",
		"email_verified": true,
	}

	userData, _ := json.Marshal(auth0User)
	req, err := http.NewRequestWithContext(ctx, "POST", 
		fmt.Sprintf("https://%s/api/v2/users", a.config.Domain), 
		strings.NewReader(string(userData)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+mgmtToken.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create Auth0 user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create Auth0 user: %d", resp.StatusCode)
	}

	// Parse the created user
	var createdUser Auth0User
	if err := json.NewDecoder(resp.Body).Decode(&createdUser); err != nil {
		return nil, fmt.Errorf("failed to decode Auth0 user: %w", err)
	}

	// Convert to our User model
	user := &User{
		ID:        createdUser.UserID,
		Email:     createdUser.Email,
		Name:      createdUser.Name,
		Picture:   createdUser.Picture,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return user, nil
}

// RegisterWithSocial registers a new user with social login
func (a *Auth0Provider) RegisterWithSocial(ctx context.Context, provider, token string) (*User, error) {
	// Validate the social token
	claims, err := a.validateJWT(token)
	if err != nil {
		return nil, fmt.Errorf("invalid social token: %w", err)
	}

	// Extract user information from claims
	user := &User{
		ID:        claims["sub"].(string),
		Email:     claims["email"].(string),
		Name:      claims["name"].(string),
		Picture:   claims["picture"].(string),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return user, nil
}

// GetUser retrieves a user by ID from Auth0
func (a *Auth0Provider) GetUser(ctx context.Context, userID string) (*User, error) {
	// Get management API token
	mgmtToken, err := a.getManagementToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get management token: %w", err)
	}

	// Get user from Auth0
	req, err := http.NewRequestWithContext(ctx, "GET", 
		fmt.Sprintf("https://%s/api/v2/users/%s", a.config.Domain, userID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+mgmtToken.AccessToken)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get Auth0 user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get Auth0 user: %d", resp.StatusCode)
	}

	// Parse the user
	var auth0User Auth0User
	if err := json.NewDecoder(resp.Body).Decode(&auth0User); err != nil {
		return nil, fmt.Errorf("failed to decode Auth0 user: %w", err)
	}

	// Convert to our User model
	user := &User{
		ID:        auth0User.UserID,
		Email:     auth0User.Email,
		Name:      auth0User.Name,
		Picture:   auth0User.Picture,
		IsActive:  true,
		UpdatedAt: time.Now(),
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email from Auth0
func (a *Auth0Provider) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	// Get management API token
	mgmtToken, err := a.getManagementToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get management token: %w", err)
	}

	// Search for user by email
	req, err := http.NewRequestWithContext(ctx, "GET", 
		fmt.Sprintf("https://%s/api/v2/users-by-email?email=%s", a.config.Domain, email), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+mgmtToken.AccessToken)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to search Auth0 user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to search Auth0 user: %d", resp.StatusCode)
	}

	// Parse the users
	var auth0Users []Auth0User
	if err := json.NewDecoder(resp.Body).Decode(&auth0Users); err != nil {
		return nil, fmt.Errorf("failed to decode Auth0 users: %w", err)
	}

	if len(auth0Users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	// Convert to our User model
	auth0User := auth0Users[0]
	user := &User{
		ID:        auth0User.UserID,
		Email:     auth0User.Email,
		Name:      auth0User.Name,
		Picture:   auth0User.Picture,
		IsActive:  true,
		UpdatedAt: time.Now(),
	}

	return user, nil
}

// UpdateUser updates user information in Auth0
func (a *Auth0Provider) UpdateUser(ctx context.Context, user *User) error {
	// Get management API token
	mgmtToken, err := a.getManagementToken()
	if err != nil {
		return fmt.Errorf("failed to get management token: %w", err)
	}

	// Update user in Auth0
	updateData := map[string]interface{}{
		"name":  user.Name,
		"email": user.Email,
	}

	userData, _ := json.Marshal(updateData)
	req, err := http.NewRequestWithContext(ctx, "PATCH", 
		fmt.Sprintf("https://%s/api/v2/users/%s", a.config.Domain, user.ID), 
		strings.NewReader(string(userData)))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+mgmtToken.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update Auth0 user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update Auth0 user: %d", resp.StatusCode)
	}

	return nil
}

// DeleteUser deletes a user from Auth0
func (a *Auth0Provider) DeleteUser(ctx context.Context, userID string) error {
	// Get management API token
	mgmtToken, err := a.getManagementToken()
	if err != nil {
		return fmt.Errorf("failed to get management token: %w", err)
	}

	// Delete user from Auth0
	req, err := http.NewRequestWithContext(ctx, "DELETE", 
		fmt.Sprintf("https://%s/api/v2/users/%s", a.config.Domain, userID), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+mgmtToken.AccessToken)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete Auth0 user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete Auth0 user: %d", resp.StatusCode)
	}

	return nil
}

// GenerateToken generates a JWT token for a user
func (a *Auth0Provider) GenerateToken(user *User) (string, error) {
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
	tokenString, err := token.SignedString([]byte(a.config.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the user
func (a *Auth0Provider) ValidateToken(tokenString string) (*User, error) {
	// Parse and validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.config.JWTSecret), nil
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
func (a *Auth0Provider) RefreshToken(refreshToken string) (string, error) {
	// For simplicity, we'll just generate a new token
	// In a real implementation, you'd validate the refresh token and generate a new access token
	return "", fmt.Errorf("refresh token not implemented")
}

// SendInvitation sends an invitation email to a user
func (a *Auth0Provider) SendInvitation(ctx context.Context, email, companyID string, invitedBy string) error {
	// Get management API token
	mgmtToken, err := a.getManagementToken()
	if err != nil {
		return fmt.Errorf("failed to get management token: %w", err)
	}

	// Create invitation in Auth0
	invitationData := map[string]interface{}{
		"client_id": a.config.ClientID,
		"email":     email,
		"connection": "Username-Password-Authentication",
		"app_metadata": map[string]interface{}{
			"company_id": companyID,
			"invited_by": invitedBy,
		},
	}

	data, _ := json.Marshal(invitationData)
	req, err := http.NewRequestWithContext(ctx, "POST", 
		fmt.Sprintf("https://%s/api/v2/jobs/verification-email", a.config.Domain), 
		strings.NewReader(string(data)))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+mgmtToken.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send invitation: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to send invitation: %d", resp.StatusCode)
	}

	return nil
}

// ActivateUser activates a user account
func (a *Auth0Provider) ActivateUser(ctx context.Context, activationToken string) error {
	// This would typically be handled by Auth0's email verification flow
	// For now, we'll just return success
	return nil
}

// ResetPassword sends a password reset email
func (a *Auth0Provider) ResetPassword(ctx context.Context, email string) error {
	// Get management API token
	mgmtToken, err := a.getManagementToken()
	if err != nil {
		return fmt.Errorf("failed to get management token: %w", err)
	}

	// Send password reset email
	resetData := map[string]interface{}{
		"client_id": a.config.ClientID,
		"email":     email,
		"connection": "Username-Password-Authentication",
	}

	data, _ := json.Marshal(resetData)
	req, err := http.NewRequestWithContext(ctx, "POST", 
		fmt.Sprintf("https://%s/api/v2/jobs/verification-email", a.config.Domain), 
		strings.NewReader(string(data)))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+mgmtToken.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send password reset: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to send password reset: %d", resp.StatusCode)
	}

	return nil
}

// ChangePassword changes a user's password
func (a *Auth0Provider) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	// Get management API token
	mgmtToken, err := a.getManagementToken()
	if err != nil {
		return fmt.Errorf("failed to get management token: %w", err)
	}

	// Update password in Auth0
	passwordData := map[string]interface{}{
		"password": newPassword,
	}

	data, _ := json.Marshal(passwordData)
	req, err := http.NewRequestWithContext(ctx, "PATCH", 
		fmt.Sprintf("https://%s/api/v2/users/%s", a.config.Domain, userID), 
		strings.NewReader(string(data)))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+mgmtToken.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to change password: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to change password: %d", resp.StatusCode)
	}

	return nil
}

// validateJWT validates an Auth0 JWT token
func (a *Auth0Provider) validateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// In a real implementation, you'd fetch the public key from Auth0
		// For now, we'll use the JWT secret
		return []byte(a.config.JWTSecret), nil
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

// getManagementToken gets an Auth0 management API token
func (a *Auth0Provider) getManagementToken() (*Auth0ManagementToken, error) {
	// Request management API token
	tokenData := map[string]string{
		"client_id":     a.config.ClientID,
		"client_secret": a.config.ClientSecret,
		"audience":      fmt.Sprintf("https://%s/api/v2/", a.config.Domain),
		"grant_type":    "client_credentials",
	}

	data, _ := json.Marshal(tokenData)
	req, err := http.NewRequest("POST", 
		fmt.Sprintf("https://%s/oauth/token", a.config.Domain), 
		strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get management token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get management token: %d", resp.StatusCode)
	}

	var token Auth0ManagementToken
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("failed to decode management token: %w", err)
	}

	return &token, nil
}

