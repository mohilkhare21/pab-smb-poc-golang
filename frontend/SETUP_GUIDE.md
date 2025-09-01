# 🚀 Frontend Setup Guide

## Prerequisites

Before starting, make sure you have the following installed:

### 1. **Node.js and npm**
```bash
# Check if Node.js is installed
node --version
npm --version

# If not installed, install using Homebrew (macOS)
brew install node

# Or download from https://nodejs.org/ (LTS version recommended)
```

### 2. **Git** (if not already installed)
```bash
git --version
```

## 🛠️ Installation Steps

### Step 1: Navigate to Frontend Directory
```bash
cd /Users/mohil/work/gitrepos/pab-smb-poc-golang/frontend
```

### Step 2: Install Dependencies
```bash
# Install all required packages
npm install

# Install additional dependencies we need
npm install lucide-react
```

### Step 3: Verify Installation
```bash
# Check if all dependencies are installed
npm list --depth=0
```

## 📦 Project Structure

```
frontend/
├── src/
│   ├── components/
│   │   └── Layout/
│   │       └── DashboardLayout.tsx    # Main dashboard layout
│   ├── contexts/
│   │   └── AuthContext.tsx            # Authentication context
│   ├── pages/
│   │   ├── Login.tsx                  # Login/Register page
│   │   ├── Dashboard.tsx              # Main dashboard
│   │   └── CompanySetup.tsx           # Company setup wizard
│   ├── services/
│   │   └── api.ts                     # API service layer
│   ├── types/
│   │   └── index.ts                   # TypeScript type definitions
│   ├── App.tsx                        # Main app component
│   └── main.tsx                       # App entry point
├── package.json                        # Dependencies and scripts
├── tsconfig.json                       # TypeScript configuration
├── vite.config.ts                      # Vite configuration
└── tailwind.config.js                  # Tailwind CSS configuration
```

## 🎯 Key Features Implemented

### 1. **Authentication System**
- Login/Register functionality
- JWT token management
- Protected routes
- Authentication context

### 2. **Dashboard Layout**
- Responsive sidebar navigation
- User profile management
- Logout functionality
- Mobile-friendly design

### 3. **Company Setup Wizard**
- Multi-step setup process
- Domain configuration
- Browser customization
- User invitation setup
- Subscription activation

### 4. **API Integration**
- Complete backend API integration
- Error handling
- Loading states
- Toast notifications

## 🚀 Running the Application

### Development Mode
```bash
npm run dev
```

The application will start at `http://localhost:5173`

### Build for Production
```bash
npm run build
```

### Preview Production Build
```bash
npm run preview
```

## 🔧 Development Commands

```bash
# Start development server
npm run dev

# Build for production
npm run build

# Run tests
npm run test

# Lint code
npm run lint

# Fix linting issues
npm run lint:fix
```

## 🌐 Backend Integration

The frontend is configured to connect to the backend at:
- **Base URL**: `http://localhost:8080/api/v1`
- **Proxy**: Configured in `package.json` for development

## 🎨 UI Components

### Built with:
- **Tailwind CSS** - Utility-first CSS framework
- **Lucide React** - Beautiful, customizable icons
- **React Hook Form** - Performant forms with easy validation
- **React Query** - Powerful data fetching and caching

### Design System:
- **Color Palette**: Blue-based primary colors
- **Typography**: Clean, readable fonts
- **Spacing**: Consistent 4px grid system
- **Components**: Modern, accessible UI components

## 📱 Responsive Design

- **Mobile-first** approach
- **Breakpoints**: sm (640px), md (768px), lg (1024px), xl (1280px)
- **Touch-friendly** interactions
- **Collapsible sidebar** for mobile

## 🔐 Authentication Flow

1. **Login/Register** → User enters credentials
2. **Token Storage** → JWT stored in localStorage
3. **Route Protection** → Unauthenticated users redirected to login
4. **Auto-logout** → Token expiration handling
5. **Persistent Sessions** → Automatic token verification

## 🚧 Next Steps

### Immediate Tasks:
1. Install dependencies: `npm install`
2. Start development server: `npm run dev`
3. Test authentication flow
4. Verify backend connectivity

### Upcoming Features:
1. **User Management** - Invite, manage, and track users
2. **Browser Shortcuts** - Custom shortcut management
3. **Configuration Panel** - Feature configuration interface
4. **Download Portal** - Browser download interface
5. **Analytics Dashboard** - Usage statistics and insights

## 🐛 Troubleshooting

### Common Issues:

#### 1. **Port Already in Use**
```bash
# Kill process on port 5173
lsof -ti:5173 | xargs kill -9
```

#### 2. **Dependencies Not Found**
```bash
# Clear npm cache and reinstall
npm cache clean --force
rm -rf node_modules package-lock.json
npm install
```

#### 3. **TypeScript Errors**
```bash
# Check TypeScript configuration
npx tsc --noEmit
```

#### 4. **Backend Connection Issues**
- Verify backend is running on port 8080
- Check CORS configuration in backend
- Verify API endpoints are accessible

## 📚 Additional Resources

- **Tailwind CSS**: https://tailwindcss.com/docs
- **React Router**: https://reactrouter.com/
- **React Query**: https://tanstack.com/query/latest
- **Lucide Icons**: https://lucide.dev/
- **Vite**: https://vitejs.dev/

## 🎉 Ready to Start!

Once you've completed the setup:

1. **Backend should be running** on port 8080
2. **Frontend will start** on port 5173
3. **Navigate to** `http://localhost:5173`
4. **Test the login flow** with your backend

The application is fully integrated with your enhanced backend and ready for development! 🚀
