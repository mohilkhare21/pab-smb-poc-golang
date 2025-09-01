# Multi-Tenant Admin Portal POC

A proof-of-concept for a multi-tenant web application built with Golang backend and React frontend.

## Architecture

### Backend (Golang)
- **Clean Architecture**: Interfaces for easy provider swapping
- **Auth Provider Interface**: Supports Auth0, Google OAuth, and custom auth
- **Database Interface**: Supports Firestore, Cloud SQL, and other databases
- **Multi-tenant Support**: Isolated data per customer/company
- **RESTful API**: JSON endpoints for frontend communication

### Frontend (React)
- **Admin Portal**: Multi-step onboarding flow
- **Management Dashboard**: User management, subscription handling
- **Download Center**: Browser distribution for different platforms
- **Responsive Design**: Works on desktop and mobile

## Features

### Onboarding Flow
1. **Admin Signup**: Social login (Google) or username/password
2. **Company Setup**: Domain configuration, color themes, user invitations
3. **Browser Shortcuts**: Auto-discovery based on company domain
4. **Payment Integration**: Stripe integration with free trial
5. **User Management**: Invite, activate, and manage team members

### Multi-tenant Features
- **Domain-based Isolation**: Each company gets isolated data
- **Custom Branding**: Color themes and company-specific configurations
- **User Roles**: Admin and regular user permissions
- **Subscription Management**: Billing and trial management

## Technology Stack

### Backend
- **Language**: Go 1.25+
- **Framework**: Gin (HTTP router)
- **Auth**: Auth0 (with interface for other providers)
- **Database**: Firestore (with interface for other databases)
- **Payment**: Stripe integration
- **External Integration**: Custom onboarding endpoints

### Frontend
- **Framework**: React 18+ with TypeScript
- **State Management**: React Context + Hooks
- **UI Library**: Material-UI or Tailwind CSS
- **HTTP Client**: Axios
- **Routing**: React Router

## Project Structure

```
├── backend/                 # Golang backend
│   ├── cmd/                # Application entry points
│   ├── internal/           # Private application code
│   │   ├── auth/          # Authentication interfaces and implementations
│   │   ├── database/      # Database interfaces and implementations
│   │   ├── handlers/      # HTTP request handlers
│   │   ├── middleware/    # HTTP middleware
│   │   ├── models/        # Data models
│   │   └── services/      # Business logic
│   ├── pkg/               # Public packages
│   └── configs/           # Configuration files
├── frontend/              # React frontend
│   ├── src/
│   │   ├── components/    # Reusable UI components
│   │   ├── pages/         # Page components
│   │   ├── services/      # API service calls
│   │   ├── hooks/         # Custom React hooks
│   │   ├── context/       # React context providers
│   │   └── types/         # TypeScript type definitions
│   └── public/            # Static assets
└── docs/                  # Documentation
```

## Getting Started

### Prerequisites
- Go 1.25+
- Node.js 18+
- Firebase project with Firestore
- Auth0 account
- Stripe account (for payments)

### Backend Setup
```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

### Frontend Setup
```bash
cd frontend
npm install
npm start
```

## Environment Variables

### Backend
- `AUTH0_DOMAIN`: Auth0 domain
- `AUTH0_CLIENT_ID`: Auth0 client ID
- `AUTH0_CLIENT_SECRET`: Auth0 client secret
- `FIRESTORE_PROJECT_ID`: Firebase project ID
- `STRIPE_SECRET_KEY`: Stripe secret key
- `JWT_SECRET`: JWT signing secret

### Frontend
- `REACT_APP_API_URL`: Backend API URL
- `REACT_APP_AUTH0_DOMAIN`: Auth0 domain
- `REACT_APP_AUTH0_CLIENT_ID`: Auth0 client ID
- `REACT_APP_STRIPE_PUBLISHABLE_KEY`: Stripe publishable key

## Development Roadmap

### Phase 1: Core Infrastructure
- [x] Project setup and structure
- [ ] Backend interfaces for auth and database
- [ ] Basic API endpoints
- [ ] Frontend project setup

### Phase 2: Authentication & User Management
- [ ] Auth0 integration
- [ ] User registration and login
- [ ] Multi-tenant user isolation
- [ ] Admin role assignment

### Phase 3: Company Onboarding
- [ ] Company domain setup
- [ ] User invitation system
- [ ] Color theme configuration
- [ ] Browser shortcut discovery

### Phase 4: Payment Integration
- [ ] Stripe integration
- [ ] Subscription management
- [ ] Free trial implementation

### Phase 5: Management Dashboard
- [ ] User management interface
- [ ] Subscription management
- [ ] Download center
- [ ] External system integration

## Contributing

This is a POC project. The architecture is designed to be easily extensible and maintainable for future development.
