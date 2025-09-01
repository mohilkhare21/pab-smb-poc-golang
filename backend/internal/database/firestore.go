package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FirestoreProvider implements DatabaseProvider for Firestore
type FirestoreProvider struct {
	client *firestore.Client
	ctx    context.Context
}

// NewFirestoreProvider creates a new Firestore provider
func NewFirestoreProvider(config DatabaseConfig) (*FirestoreProvider, error) {
	ctx := context.Background()
	
	var opts []option.ClientOption
	if config.Credentials != "" {
		opts = append(opts, option.WithCredentialsFile(config.Credentials))
	}
	
	client, err := firestore.NewClient(ctx, config.ProjectID, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create firestore client: %w", err)
	}
	
	return &FirestoreProvider{
		client: client,
		ctx:    ctx,
	}, nil
}

// CreateCompany creates a new company
func (f *FirestoreProvider) CreateCompany(ctx context.Context, company *Company) error {
	company.CreatedAt = time.Now()
	company.UpdatedAt = time.Now()
	company.Status = "trial"
	company.TrialEndsAt = time.Now().AddDate(0, 1, 0) // 1 month trial
	
	_, err := f.client.Collection("companies").Doc(company.ID).Set(ctx, company)
	return err
}

// GetCompany retrieves a company by ID
func (f *FirestoreProvider) GetCompany(ctx context.Context, companyID string) (*Company, error) {
	doc, err := f.client.Collection("companies").Doc(companyID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("company not found")
		}
		return nil, err
	}
	
	var company Company
	if err := doc.DataTo(&company); err != nil {
		return nil, err
	}
	
	return &company, nil
}

// GetCompanyByDomain retrieves a company by domain
func (f *FirestoreProvider) GetCompanyByDomain(ctx context.Context, domain string) (*Company, error) {
	iter := f.client.Collection("companies").Where("domain", "==", domain).Limit(1).Documents(ctx)
	defer iter.Stop()
	
	doc, err := iter.Next()
	if err != nil {
		return nil, fmt.Errorf("company not found")
	}
	
	var company Company
	if err := doc.DataTo(&company); err != nil {
		return nil, err
	}
	
	return &company, nil
}

// UpdateCompany updates a company
func (f *FirestoreProvider) UpdateCompany(ctx context.Context, company *Company) error {
	company.UpdatedAt = time.Now()
	
	_, err := f.client.Collection("companies").Doc(company.ID).Set(ctx, company)
	return err
}

// DeleteCompany deletes a company
func (f *FirestoreProvider) DeleteCompany(ctx context.Context, companyID string) error {
	_, err := f.client.Collection("companies").Doc(companyID).Delete(ctx)
	return err
}

// ListCompanies lists companies with pagination
func (f *FirestoreProvider) ListCompanies(ctx context.Context, limit, offset int) ([]*Company, error) {
	iter := f.client.Collection("companies").OrderBy("created_at", firestore.Desc).Limit(limit).Offset(offset).Documents(ctx)
	defer iter.Stop()
	
	var companies []*Company
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		
		var company Company
		if err := doc.DataTo(&company); err != nil {
			continue
		}
		companies = append(companies, &company)
	}
	
	return companies, nil
}

// CreateUser creates a new user
func (f *FirestoreProvider) CreateUser(ctx context.Context, user *User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsActive = true
	
	_, err := f.client.Collection("users").Doc(user.ID).Set(ctx, user)
	return err
}

// GetUser retrieves a user by ID
func (f *FirestoreProvider) GetUser(ctx context.Context, userID string) (*User, error) {
	doc, err := f.client.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	
	var user User
	if err := doc.DataTo(&user); err != nil {
		return nil, err
	}
	
	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (f *FirestoreProvider) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	iter := f.client.Collection("users").Where("email", "==", email).Limit(1).Documents(ctx)
	defer iter.Stop()
	
	doc, err := iter.Next()
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	
	var user User
	if err := doc.DataTo(&user); err != nil {
		return nil, err
	}
	
	return &user, nil
}

// GetUsersByCompany retrieves all users for a company
func (f *FirestoreProvider) GetUsersByCompany(ctx context.Context, companyID string) ([]*User, error) {
	iter := f.client.Collection("users").Where("company_id", "==", companyID).Documents(ctx)
	defer iter.Stop()
	
	var users []*User
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		
		var user User
		if err := doc.DataTo(&user); err != nil {
			continue
		}
		users = append(users, &user)
	}
	
	return users, nil
}

// UpdateUser updates a user
func (f *FirestoreProvider) UpdateUser(ctx context.Context, user *User) error {
	user.UpdatedAt = time.Now()
	
	_, err := f.client.Collection("users").Doc(user.ID).Set(ctx, user)
	return err
}

// DeleteUser deletes a user
func (f *FirestoreProvider) DeleteUser(ctx context.Context, userID string) error {
	_, err := f.client.Collection("users").Doc(userID).Delete(ctx)
	return err
}

// CountUsersByCompany counts users in a company
func (f *FirestoreProvider) CountUsersByCompany(ctx context.Context, companyID string) (int, error) {
	iter := f.client.Collection("users").Where("company_id", "==", companyID).Documents(ctx)
	defer iter.Stop()
	
	count := 0
	for {
		_, err := iter.Next()
		if err != nil {
			break
		}
		count++
	}
	
	return count, nil
}

// CreateInvitation creates a new invitation
func (f *FirestoreProvider) CreateInvitation(ctx context.Context, invitation *Invitation) error {
	invitation.CreatedAt = time.Now()
	invitation.Status = "pending"
	invitation.ExpiresAt = time.Now().AddDate(0, 0, 7) // 7 days expiry
	
	_, err := f.client.Collection("invitations").Doc(invitation.ID).Set(ctx, invitation)
	return err
}

// GetInvitation retrieves an invitation by ID
func (f *FirestoreProvider) GetInvitation(ctx context.Context, invitationID string) (*Invitation, error) {
	doc, err := f.client.Collection("invitations").Doc(invitationID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("invitation not found")
		}
		return nil, err
	}
	
	var invitation Invitation
	if err := doc.DataTo(&invitation); err != nil {
		return nil, err
	}
	
	return &invitation, nil
}

// GetInvitationByToken retrieves an invitation by token
func (f *FirestoreProvider) GetInvitationByToken(ctx context.Context, token string) (*Invitation, error) {
	iter := f.client.Collection("invitations").Where("token", "==", token).Limit(1).Documents(ctx)
	defer iter.Stop()
	
	doc, err := iter.Next()
	if err != nil {
		return nil, fmt.Errorf("invitation not found")
	}
	
	var invitation Invitation
	if err := doc.DataTo(&invitation); err != nil {
		return nil, err
	}
	
	return &invitation, nil
}

// GetInvitationsByCompany retrieves all invitations for a company
func (f *FirestoreProvider) GetInvitationsByCompany(ctx context.Context, companyID string) ([]*Invitation, error) {
	iter := f.client.Collection("invitations").Where("company_id", "==", companyID).Documents(ctx)
	defer iter.Stop()
	
	var invitations []*Invitation
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		
		var invitation Invitation
		if err := doc.DataTo(&invitation); err != nil {
			continue
		}
		invitations = append(invitations, &invitation)
	}
	
	return invitations, nil
}

// UpdateInvitation updates an invitation
func (f *FirestoreProvider) UpdateInvitation(ctx context.Context, invitation *Invitation) error {
	_, err := f.client.Collection("invitations").Doc(invitation.ID).Set(ctx, invitation)
	return err
}

// DeleteInvitation deletes an invitation
func (f *FirestoreProvider) DeleteInvitation(ctx context.Context, invitationID string) error {
	_, err := f.client.Collection("invitations").Doc(invitationID).Delete(ctx)
	return err
}

// DeleteExpiredInvitations deletes expired invitations
func (f *FirestoreProvider) DeleteExpiredInvitations(ctx context.Context) error {
	iter := f.client.Collection("invitations").Where("expires_at", "<", time.Now()).Documents(ctx)
	defer iter.Stop()
	
	batch := f.client.Batch()
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		batch.Delete(doc.Ref)
	}
	
	_, err := batch.Commit(ctx)
	return err
}

// CreateBrowserShortcut creates a new browser shortcut
func (f *FirestoreProvider) CreateBrowserShortcut(ctx context.Context, shortcut *BrowserShortcut) error {
	_, err := f.client.Collection("browser_shortcuts").Doc(shortcut.ID).Set(ctx, shortcut)
	return err
}

// GetBrowserShortcut retrieves a browser shortcut by ID
func (f *FirestoreProvider) GetBrowserShortcut(ctx context.Context, shortcutID string) (*BrowserShortcut, error) {
	doc, err := f.client.Collection("browser_shortcuts").Doc(shortcutID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("browser shortcut not found")
		}
		return nil, err
	}
	
	var shortcut BrowserShortcut
	if err := doc.DataTo(&shortcut); err != nil {
		return nil, err
	}
	
	return &shortcut, nil
}

// GetBrowserShortcutsByCompany retrieves all browser shortcuts for a company
func (f *FirestoreProvider) GetBrowserShortcutsByCompany(ctx context.Context, companyID string) ([]*BrowserShortcut, error) {
	iter := f.client.Collection("browser_shortcuts").Where("company_id", "==", companyID).OrderBy("order", firestore.Asc).Documents(ctx)
	defer iter.Stop()
	
	var shortcuts []*BrowserShortcut
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		
		var shortcut BrowserShortcut
		if err := doc.DataTo(&shortcut); err != nil {
			continue
		}
		shortcuts = append(shortcuts, &shortcut)
	}
	
	return shortcuts, nil
}

// UpdateBrowserShortcut updates a browser shortcut
func (f *FirestoreProvider) UpdateBrowserShortcut(ctx context.Context, shortcut *BrowserShortcut) error {
	_, err := f.client.Collection("browser_shortcuts").Doc(shortcut.ID).Set(ctx, shortcut)
	return err
}

// DeleteBrowserShortcut deletes a browser shortcut
func (f *FirestoreProvider) DeleteBrowserShortcut(ctx context.Context, shortcutID string) error {
	_, err := f.client.Collection("browser_shortcuts").Doc(shortcutID).Delete(ctx)
	return err
}

// DeleteBrowserShortcutsByCompany deletes all browser shortcuts for a company
func (f *FirestoreProvider) DeleteBrowserShortcutsByCompany(ctx context.Context, companyID string) error {
	iter := f.client.Collection("browser_shortcuts").Where("company_id", "==", companyID).Documents(ctx)
	defer iter.Stop()
	
	batch := f.client.Batch()
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		batch.Delete(doc.Ref)
	}
	
	_, err := batch.Commit(ctx)
	return err
}

// CreateSubscription creates a new subscription
func (f *FirestoreProvider) CreateSubscription(ctx context.Context, subscription *Subscription) error {
	subscription.CreatedAt = time.Now()
	subscription.UpdatedAt = time.Now()
	
	_, err := f.client.Collection("subscriptions").Doc(subscription.ID).Set(ctx, subscription)
	return err
}

// GetSubscription retrieves a subscription by ID
func (f *FirestoreProvider) GetSubscription(ctx context.Context, subscriptionID string) (*Subscription, error) {
	doc, err := f.client.Collection("subscriptions").Doc(subscriptionID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("subscription not found")
		}
		return nil, err
	}
	
	var subscription Subscription
	if err := doc.DataTo(&subscription); err != nil {
		return nil, err
	}
	
	return &subscription, nil
}

// GetSubscriptionByCompany retrieves a subscription by company ID
func (f *FirestoreProvider) GetSubscriptionByCompany(ctx context.Context, companyID string) (*Subscription, error) {
	iter := f.client.Collection("subscriptions").Where("company_id", "==", companyID).Limit(1).Documents(ctx)
	defer iter.Stop()
	
	doc, err := iter.Next()
	if err != nil {
		return nil, fmt.Errorf("subscription not found")
	}
	
	var subscription Subscription
	if err := doc.DataTo(&subscription); err != nil {
		return nil, err
	}
	
	return &subscription, nil
}

// UpdateSubscription updates a subscription
func (f *FirestoreProvider) UpdateSubscription(ctx context.Context, subscription *Subscription) error {
	subscription.UpdatedAt = time.Now()
	
	_, err := f.client.Collection("subscriptions").Doc(subscription.ID).Set(ctx, subscription)
	return err
}

// DeleteSubscription deletes a subscription
func (f *FirestoreProvider) DeleteSubscription(ctx context.Context, subscriptionID string) error {
	_, err := f.client.Collection("subscriptions").Doc(subscriptionID).Delete(ctx)
	return err
}

// BeginTransaction starts a new transaction
func (f *FirestoreProvider) BeginTransaction(ctx context.Context) (Transaction, error) {
	// For Firestore, we'll use a batch operation as a simple transaction
	batch := f.client.Batch()
	
	return &FirestoreTransaction{batch: batch}, nil
}

// Ping checks if the database is accessible
func (f *FirestoreProvider) Ping(ctx context.Context) error {
	// For Firestore, we just need to verify the client is working
	// Try to list a collection to verify connectivity
	iter := f.client.Collection("_health").Limit(1).Documents(ctx)
	defer iter.Stop()
	
	// Just check if we can make the request, don't care about results
	_, err := iter.Next()
	if err != nil && err != iterator.Done {
		return err
	}
	
	return nil
}

// Close closes the database connection
func (f *FirestoreProvider) Close() error {
	return f.client.Close()
}

// FirestoreTransaction implements Transaction for Firestore
type FirestoreTransaction struct {
	batch *firestore.WriteBatch
}

// Commit commits the transaction
func (t *FirestoreTransaction) Commit() error {
	// Firestore batch needs to be committed
	_, err := t.batch.Commit(context.Background())
	return err
}

// Rollback rolls back the transaction
func (t *FirestoreTransaction) Rollback() error {
	// Firestore doesn't support rollback for batches
	log.Println("Warning: Firestore doesn't support rollback for batches")
	return nil
}

// Enhanced User Operations

// GetInvitedUsersByCompany gets users with invitation status
func (f *FirestoreProvider) GetInvitedUsersByCompany(ctx context.Context, companyID string) ([]*User, error) {
	iter := f.client.Collection("users").Where("company_id", "==", companyID).Where("invitation_status", "==", "invited").Documents(ctx)
	defer iter.Stop()
	
	var users []*User
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		
		var user User
		if err := doc.DataTo(&user); err != nil {
			continue
		}
		users = append(users, &user)
	}
	
	return users, nil
}

// GetActiveUsersByCompany gets active users
func (f *FirestoreProvider) GetActiveUsersByCompany(ctx context.Context, companyID string) ([]*User, error) {
	iter := f.client.Collection("users").Where("company_id", "==", companyID).Where("is_active", "==", true).Documents(ctx)
	defer iter.Stop()
	
	var users []*User
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		
		var user User
		if err := doc.DataTo(&user); err != nil {
			continue
		}
		users = append(users, &user)
	}
	
	return users, nil
}

// CountInvitedUsersByCompany counts invited users
func (f *FirestoreProvider) CountInvitedUsersByCompany(ctx context.Context, companyID string) (int, error) {
	iter := f.client.Collection("users").Where("company_id", "==", companyID).Where("invitation_status", "==", "invited").Documents(ctx)
	defer iter.Stop()
	
	count := 0
	for {
		_, err := iter.Next()
		if err != nil {
			break
		}
		count++
	}
	
	return count, nil
}

// CountActiveUsersByCompany counts active users
func (f *FirestoreProvider) CountActiveUsersByCompany(ctx context.Context, companyID string) (int, error) {
	iter := f.client.Collection("users").Where("company_id", "==", companyID).Where("is_active", "==", true).Documents(ctx)
	defer iter.Stop()
	
	count := 0
	for {
		_, err := iter.Next()
		if err != nil {
			break
		}
		count++
	}
	
	return count, nil
}

// UpdateUserInvitationStatus updates user invitation status
func (f *FirestoreProvider) UpdateUserInvitationStatus(ctx context.Context, userID string, status string) error {
	user, err := f.GetUser(ctx, userID)
	if err != nil {
		return err
	}
	
	user.InvitationStatus = status
	user.UpdatedAt = time.Now()
	
	if status == "active" {
		user.ActivatedAt = time.Now()
	}
	
	return f.UpdateUser(ctx, user)
}

// Enhanced Invitation Operations

// GetPendingInvitationsByCompany gets pending invitations
func (f *FirestoreProvider) GetPendingInvitationsByCompany(ctx context.Context, companyID string) ([]*Invitation, error) {
	iter := f.client.Collection("invitations").Where("company_id", "==", companyID).Where("status", "==", "pending").Documents(ctx)
	defer iter.Stop()
	
	var invitations []*Invitation
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		
		var invitation Invitation
		if err := doc.DataTo(&invitation); err != nil {
			continue
		}
		invitations = append(invitations, &invitation)
	}
	
	return invitations, nil
}

// CountPendingInvitationsByCompany counts pending invitations
func (f *FirestoreProvider) CountPendingInvitationsByCompany(ctx context.Context, companyID string) (int, error) {
	iter := f.client.Collection("invitations").Where("company_id", "==", companyID).Where("status", "==", "pending").Documents(ctx)
	defer iter.Stop()
	
	count := 0
	for {
		_, err := iter.Next()
		if err != nil {
			break
		}
		count++
	}
	
	return count, nil
}

// UpdateInvitationSentStatus updates invitation sent status
func (f *FirestoreProvider) UpdateInvitationSentStatus(ctx context.Context, invitationID string, sentAt time.Time) error {
	invitation, err := f.GetInvitation(ctx, invitationID)
	if err != nil {
		return err
	}
	
	invitation.Status = "sent"
	invitation.SentAt = sentAt
	invitation.SentCount++
	invitation.LastSentAt = sentAt
	
	return f.UpdateInvitation(ctx, invitation)
}

// ResendInvitation resends an invitation
func (f *FirestoreProvider) ResendInvitation(ctx context.Context, invitationID string) error {
	invitation, err := f.GetInvitation(ctx, invitationID)
	if err != nil {
		return err
	}
	
	// Update sent status
	now := time.Now()
	invitation.SentAt = now
	invitation.SentCount++
	invitation.LastSentAt = now
	
	return f.UpdateInvitation(ctx, invitation)
}

// Enhanced Shortcut Operations

// GetSuggestedShortcutsByCompany gets suggested shortcuts
func (f *FirestoreProvider) GetSuggestedShortcutsByCompany(ctx context.Context, companyID string) ([]*BrowserShortcut, error) {
	iter := f.client.Collection("browser_shortcuts").Where("company_id", "==", companyID).Where("is_suggested", "==", true).Documents(ctx)
	defer iter.Stop()
	
	var shortcuts []*BrowserShortcut
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		
		var shortcut BrowserShortcut
		if err := doc.DataTo(&shortcut); err != nil {
			continue
		}
		shortcuts = append(shortcuts, &shortcut)
	}
	
	return shortcuts, nil
}

// GetCustomShortcutsByCompany gets custom shortcuts
func (f *FirestoreProvider) GetCustomShortcutsByCompany(ctx context.Context, companyID string) ([]*BrowserShortcut, error) {
	iter := f.client.Collection("browser_shortcuts").Where("company_id", "==", companyID).Where("category", "==", "custom").Documents(ctx)
	defer iter.Stop()
	
	var shortcuts []*BrowserShortcut
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		
		var shortcut BrowserShortcut
		if err := doc.DataTo(&shortcut); err != nil {
			continue
		}
		shortcuts = append(shortcuts, &shortcut)
	}
	
	return shortcuts, nil
}

// GenerateShortcutsForDomain generates suggested shortcuts for a domain
func (f *FirestoreProvider) GenerateShortcutsForDomain(ctx context.Context, companyID string, domain string) error {
	// Common shortcuts based on domain analysis
	commonShortcuts := []BrowserShortcut{
		{
			ID:          fmt.Sprintf("shortcut_%s_gmail", companyID),
			CompanyID:   companyID,
			Name:        "Gmail",
			URL:         "https://mail.google.com",
			Icon:        "gmail-icon",
			Description: "Access Gmail",
			Order:       1,
			IsActive:    true,
			IsSuggested: true,
			Category:    "suggested",
			Source:      "auto-generated",
		},
		{
			ID:          fmt.Sprintf("shortcut_%s_calendar", companyID),
			CompanyID:   companyID,
			Name:        "Google Calendar",
			URL:         "https://calendar.google.com",
			Icon:        "calendar-icon",
			Description: "Access Google Calendar",
			Order:       2,
			IsActive:    true,
			IsSuggested: true,
			Category:    "suggested",
			Source:      "auto-generated",
		},
		{
			ID:          fmt.Sprintf("shortcut_%s_drive", companyID),
			CompanyID:   companyID,
			Name:        "Google Drive",
			URL:         "https://drive.google.com",
			Icon:        "drive-icon",
			Description: "Access Google Drive",
			Order:       3,
			IsActive:    true,
			IsSuggested: true,
			Category:    "suggested",
			Source:      "auto-generated",
		},
		{
			ID:          fmt.Sprintf("shortcut_%s_company", companyID),
			CompanyID:   companyID,
			Name:        fmt.Sprintf("%s Website", domain),
			URL:         fmt.Sprintf("https://%s", domain),
			Icon:        "company-icon",
			Description: fmt.Sprintf("Access %s website", domain),
			Order:       0,
			IsActive:    true,
			IsSuggested: true,
			Category:    "company",
			Source:      "auto-generated",
		},
	}
	
	// Create shortcuts in batch
	batch := f.client.Batch()
	for _, shortcut := range commonShortcuts {
		docRef := f.client.Collection("browser_shortcuts").Doc(shortcut.ID)
		batch.Set(docRef, shortcut)
	}
	
	_, err := batch.Commit(ctx)
	return err
}

// Enhanced Subscription Operations

// UpdateSubscriptionUserCounts updates subscription user counts
func (f *FirestoreProvider) UpdateSubscriptionUserCounts(ctx context.Context, subscriptionID string, activeUsers, invitedUsers int) error {
	subscription, err := f.GetSubscription(ctx, subscriptionID)
	if err != nil {
		return err
	}
	
	subscription.ActiveUsers = activeUsers
	subscription.InvitedUsers = invitedUsers
	subscription.UpdatedAt = time.Now()
	
	return f.UpdateSubscription(ctx, subscription)
}

// GetSubscriptionStats gets subscription statistics
func (f *FirestoreProvider) GetSubscriptionStats(ctx context.Context, companyID string) (*Subscription, error) {
	return f.GetSubscriptionByCompany(ctx, companyID)
}

// Setup Progress Operations

// CreateSetupProgress creates setup progress
func (f *FirestoreProvider) CreateSetupProgress(ctx context.Context, progress *CompanySetupProgress) error {
	progress.LastUpdated = time.Now()
	
	_, err := f.client.Collection("setup_progress").Doc(progress.CompanyID).Set(ctx, progress)
	return err
}

// GetSetupProgress gets setup progress
func (f *FirestoreProvider) GetSetupProgress(ctx context.Context, companyID string) (*CompanySetupProgress, error) {
	doc, err := f.client.Collection("setup_progress").Doc(companyID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			// Return default progress
			return &CompanySetupProgress{
				CompanyID:    companyID,
				Step:         "domain",
				Progress:     0,
				LastUpdated:  time.Now(),
			}, nil
		}
		return nil, err
	}
	
	var progress CompanySetupProgress
	if err := doc.DataTo(&progress); err != nil {
		return nil, err
	}
	
	return &progress, nil
}

// UpdateSetupProgress updates setup progress
func (f *FirestoreProvider) UpdateSetupProgress(ctx context.Context, progress *CompanySetupProgress) error {
	progress.LastUpdated = time.Now()
	
	_, err := f.client.Collection("setup_progress").Doc(progress.CompanyID).Set(ctx, progress)
	return err
}

// UpdateSetupStep updates setup step
func (f *FirestoreProvider) UpdateSetupStep(ctx context.Context, companyID string, step string, progress int) error {
	setupProgress, err := f.GetSetupProgress(ctx, companyID)
	if err != nil {
		return err
	}
	
	setupProgress.Step = step
	setupProgress.Progress = progress
	setupProgress.LastUpdated = time.Now()
	
	return f.UpdateSetupProgress(ctx, setupProgress)
}

// Configuration Status Operations

// UpdateCompanyConfigurationStatus updates company configuration status
func (f *FirestoreProvider) UpdateCompanyConfigurationStatus(ctx context.Context, companyID string, feature string, status bool) error {
	company, err := f.GetCompany(ctx, companyID)
	if err != nil {
		return err
	}
	
	// Update the specific feature status
	switch feature {
	case "website_security":
		company.WebsiteSecurityConfigured = status
	case "malware_security":
		company.MalwareSecurityConfigured = status
	case "data_controls":
		company.DataControlsConfigured = status
	case "reporting":
		company.ReportingConfigured = status
	case "browser_customization":
		company.BrowserCustomized = status
	case "subscription":
		company.SubscriptionActive = status
	case "users_invited":
		company.UsersInvited = status
	case "download_ready":
		company.DownloadReady = status
	}
	
	company.UpdatedAt = time.Now()
	
	return f.UpdateCompany(ctx, company)
}

// GetCompanyConfigurationStatus gets company configuration status
func (f *FirestoreProvider) GetCompanyConfigurationStatus(ctx context.Context, companyID string) (map[string]bool, error) {
	company, err := f.GetCompany(ctx, companyID)
	if err != nil {
		return nil, err
	}
	
	status := map[string]bool{
		"website_security":    company.WebsiteSecurityConfigured,
		"malware_security":    company.MalwareSecurityConfigured,
		"data_controls":       company.DataControlsConfigured,
		"reporting":           company.ReportingConfigured,
		"browser_customization": company.BrowserCustomized,
		"subscription":        company.SubscriptionActive,
		"users_invited":       company.UsersInvited,
		"download_ready":      company.DownloadReady,
	}
	
	return status, nil
}
