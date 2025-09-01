# Multi-Tenant Admin Portal API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
Most endpoints require authentication using JWT tokens. Include the token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

## Endpoints

### Authentication

#### POST /auth/login
Login with email and password.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "jwt-token-here",
    "user": {
      "id": "user-id",
      "email": "user@example.com",
      "name": "John Doe",
      "picture": "https://example.com/avatar.jpg",
      "company_id": "company-id",
      "role": "admin"
    }
  }
}
```

#### POST /auth/register
Register a new user.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Registration successful",
  "data": {
    "token": "jwt-token-here",
    "user": {
      "id": "user-id",
      "email": "user@example.com",
      "name": "John Doe",
      "picture": "https://example.com/avatar.jpg",
      "role": "admin"
    }
  }
}
```

#### GET /auth/verify
Verify JWT token and get user information.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response:**
```json
{
  "success": true,
  "message": "Authentication successful",
  "data": {
    "user": {
      "id": "user-id",
      "email": "user@example.com",
      "name": "John Doe",
      "picture": "https://example.com/avatar.jpg",
      "company_id": "company-id",
      "role": "admin",
      "is_active": true
    }
  }
}
```

#### POST /auth/reset-password
Request password reset.

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

### Company Management

#### POST /companies
Create a new company (requires authentication).

**Request Body:**
```json
{
  "name": "Acme Corp",
  "domain": "acme.com",
  "color_theme": "#007bff"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Company created successfully",
  "data": {
    "company": {
      "id": "company-id",
      "name": "Acme Corp",
      "domain": "acme.com",
      "color_theme": "#007bff",
      "status": "trial",
      "admin_user_id": "user-id"
    }
  }
}
```

#### GET /companies/me
Get current user's company information.

**Response:**
```json
{
  "success": true,
  "data": {
    "company": {
      "id": "company-id",
      "name": "Acme Corp",
      "domain": "acme.com",
      "color_theme": "#007bff",
      "logo_url": "https://example.com/logo.png",
      "status": "trial",
      "trial_ends_at": "2024-02-01T00:00:00Z",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "onboarded": false
    }
  }
}
```

#### PUT /companies/me
Update company information (admin only).

**Request Body:**
```json
{
  "name": "Acme Corporation",
  "color_theme": "#28a745",
  "logo_url": "https://example.com/new-logo.png"
}
```

#### GET /companies/stats
Get company statistics.

**Response:**
```json
{
  "success": true,
  "data": {
    "stats": {
      "total_users": 15,
      "company_status": "trial",
      "trial_ends_at": "2024-02-01T00:00:00Z",
      "onboarded": false
    }
  }
}
```

### User Management

#### GET /users
Get all users in the company.

**Response:**
```json
{
  "success": true,
  "data": {
    "users": [
      {
        "id": "user-id",
        "email": "user@example.com",
        "name": "John Doe",
        "picture": "https://example.com/avatar.jpg",
        "role": "admin",
        "is_active": true,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z",
        "last_login_at": "2024-01-01T12:00:00Z",
        "onboarded": true
      }
    ]
  }
}
```

#### GET /users/:id
Get specific user information.

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "user-id",
      "email": "user@example.com",
      "name": "John Doe",
      "picture": "https://example.com/avatar.jpg",
      "role": "admin",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "last_login_at": "2024-01-01T12:00:00Z",
      "onboarded": true
    }
  }
}
```

#### PUT /users/:id
Update user information (admin only or self).

**Request Body:**
```json
{
  "name": "John Smith",
  "role": "user",
  "is_active": true
}
```

#### DELETE /users/:id
Delete user (admin only).

### Invitation Management

#### POST /invitations
Create user invitations (admin only).

**Request Body:**
```json
{
  "emails": ["user1@example.com", "user2@example.com"]
}
```

**Response:**
```json
{
  "success": true,
  "message": "Invitations created successfully",
  "data": {
    "invitations": [
      {
        "id": "invitation-id",
        "email": "user1@example.com",
        "status": "pending",
        "expires_at": "2024-01-08T00:00:00Z"
      }
    ]
  }
}
```

#### GET /invitations
Get all invitations for the company (admin only).

**Response:**
```json
{
  "success": true,
  "data": {
    "invitations": [
      {
        "id": "invitation-id",
        "email": "user@example.com",
        "invited_by": "admin-user-id",
        "status": "pending",
        "expires_at": "2024-01-08T00:00:00Z",
        "created_at": "2024-01-01T00:00:00Z",
        "accepted_at": null
      }
    ]
  }
}
```

#### POST /invitations/:token/accept
Accept an invitation.

**Response:**
```json
{
  "success": true,
  "message": "Invitation accepted successfully",
  "data": {
    "company_id": "company-id"
  }
}
```

### Browser Shortcuts

#### GET /shortcuts
Get all browser shortcuts for the company.

**Response:**
```json
{
  "success": true,
  "data": {
    "shortcuts": [
      {
        "id": "shortcut-id",
        "name": "Gmail",
        "url": "https://mail.google.com",
        "icon": "https://example.com/gmail-icon.png",
        "description": "Access Gmail",
        "order": 1,
        "is_active": true
      }
    ]
  }
}
```

#### POST /shortcuts
Create a new browser shortcut (admin only).

**Request Body:**
```json
{
  "name": "Gmail",
  "url": "https://mail.google.com",
  "icon": "https://example.com/gmail-icon.png",
  "description": "Access Gmail",
  "order": 1
}
```

#### PUT /shortcuts/:id
Update browser shortcut (admin only).

#### DELETE /shortcuts/:id
Delete browser shortcut (admin only).

### Setup and Configuration

#### GET /setup/progress
Get the setup progress for the company.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "company_id": "company-id",
    "step": "invitations",
    "progress": 60,
    "domain_provided": true,
    "customization_completed": true,
    "invitations_sent": true,
    "subscription_started": false,
    "setup_completed": false,
    "last_updated": "2024-01-15T10:30:00Z"
  }
}
```

#### PUT /setup/step
Update the setup step for the company.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Request Body:**
```json
{
  "step": "subscription",
  "progress": 80
}
```

**Response:**
```json
{
  "success": true,
  "message": "Setup step updated successfully"
}
```

#### GET /setup/stats
Get company statistics and configuration status.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "total_users": 15,
    "active_users": 12,
    "invited_users": 3,
    "pending_invitations": 2,
    "max_users": 20,
    "setup_progress": 80,
    "configuration_status": {
      "website_security": true,
      "malware_security": true,
      "data_controls": false,
      "reporting": false,
      "browser_customization": true,
      "subscription": true,
      "users_invited": true,
      "download_ready": false
    }
  }
}
```

#### PUT /setup/config
Update configuration status for a feature.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Request Body:**
```json
{
  "feature": "data_controls",
  "status": true
}
```

**Response:**
```json
{
  "success": true,
  "message": "Configuration status updated successfully"
}
```

#### POST /setup/generate-shortcuts
Generate suggested shortcuts for a company domain.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Request Body:**
```json
{
  "domain": "example.com"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Shortcuts generated successfully"
}
```

#### POST /setup/nudge-users
Send reminders to invited users.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Request Body:**
```json
{
  "user_ids": ["user-id-1", "user-id-2"],
  "message": "Please complete your registration"
}
```

**Response:**
```json
{
  "success": true,
  "message": "User nudges sent successfully"
}
```

#### GET /setup/download-info
Get download information for the custom browser.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "download_url": "https://download.pab-smb.com/browser/latest",
    "version": "1.0.0",
    "release_date": "2024-01-15",
    "file_size": "45.2 MB",
    "supported_os": ["macOS", "Windows", "Linux"],
    "installation_instructions": "Download and run the installer. Follow the setup wizard to complete installation."
  }
}
```

## Error Responses

All endpoints return consistent error responses:

```json
{
  "success": false,
  "error": "Error message describing what went wrong"
}
```

Common HTTP status codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `409` - Conflict
- `500` - Internal Server Error

## Rate Limiting

API requests are rate limited to prevent abuse. Limits are:
- 100 requests per minute per IP address
- 1000 requests per hour per authenticated user

## Pagination

For endpoints that return lists, pagination is supported using query parameters:
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10, max: 100)

Example:
```
GET /api/v1/users?page=2&limit=20
```

## Health Check

#### GET /health
Check if the API is running.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T12:00:00Z",
  "service": "multi-tenant-admin-portal"
}
```

