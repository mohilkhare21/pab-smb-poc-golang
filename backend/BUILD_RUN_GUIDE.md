# 🚀 Build and Run Guide

## 📋 Prerequisites

Before building/running, ensure you have:

1. **✅ Go installed** - `go version` should work
2. **✅ .env file configured** - With Firebase credentials
3. **✅ Firebase service account JSON** - In `configs/` directory
4. **✅ Dependencies ready** - `go mod tidy` completed

## 🔧 **Method 1: Development Mode (Recommended)**

### **Quick Start**
```bash
# Navigate to backend directory
cd backend

# Run directly (no build step)
go run cmd/server/main.go
```

### **Expected Output**
```
2024/01/01 12:00:00 No .env file found, using system environment variables
2024/01/01 12:00:00 Database connection established
2024/01/01 12:00:00 Server starting on port 8080
2024/01/01 12:00:00 Health check available at http://localhost:8080/health
2024/01/01 12:00:00 API documentation available at http://localhost:8080/api/v1
```

## 🏗️ **Method 2: Build and Run (Production)**

### **Step 1: Build the Application**
```bash
# Build for current platform
go build -o server cmd/server/main.go

# Build for specific platform (e.g., Linux)
GOOS=linux GOARCH=amd64 go build -o server-linux cmd/server/main.go

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o server-macos cmd/server/main.go

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o server.exe cmd/server/main.go
```

### **Step 2: Run the Built Executable**
```bash
# Run the built server
./server

# Or on Windows
server.exe
```

## 🐳 **Method 3: Docker (Alternative)**

### **Build Docker Image**
```bash
# Build the image
docker build -t admin-portal-backend .

# Run the container
docker run -p 8080:8080 --env-file .env admin-portal-backend
```

## 🔍 **Method 4: Using Makefile**

### **Available Commands**
```bash
# Check available commands
make help

# Common commands
make build      # Build the application
make run        # Run the application
make clean      # Clean build artifacts
make test       # Run tests
```

## 🧪 **Testing the Build**

### **1. Health Check**
```bash
# Test if server is running
curl http://localhost:8080/health

# Expected response:
{
  "status": "healthy",
  "timestamp": "2024-01-01T12:00:00Z",
  "service": "multi-tenant-admin-portal"
}
```

### **2. API Endpoints**
```bash
# Test API base
curl http://localhost:8080/api/v1

# Test authentication endpoint
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "testpassword123",
    "name": "Test User"
  }'
```

## 🐛 **Common Build/Run Issues**

### **Issue 1: "No .env file found"**
```bash
# Solution: Ensure .env file exists in backend root
ls -la .env

# If missing, create from template
cp env.example .env
# Then edit with your actual values
```

### **Issue 2: "Database connection failed"**
```bash
# Check .env variables
cat .env | grep -E "(FIRESTORE|DB_|GOOGLE)"

# Verify Firebase credentials file exists
ls -la configs/firebase-service-account.json
```

### **Issue 3: "Port already in use"**
```bash
# Find process using port 8080
lsof -ti:8080

# Kill the process
kill -9 $(lsof -ti:8080)
```

### **Issue 4: "Module not found"**
```bash
# Clean and reinstall modules
go clean -modcache
go mod tidy
go mod download
```

## 📊 **Build Information**

### **Binary Size**
```bash
# Check binary size after build
ls -lh server

# Expected: ~10-20MB for typical Go application
```

### **Build Time**
```bash
# Time the build process
time go build -o server cmd/server/main.go

# First build: ~5-15 seconds
# Subsequent builds: ~1-3 seconds
```

## 🚀 **Performance Tips**

### **1. Fast Development Builds**
```bash
# Use go run for development (faster iteration)
go run cmd/server/main.go

# Use go build for production testing
go build -o server cmd/server/main.go
```

### **2. Optimized Production Builds**
```bash
# Build with optimizations
go build -ldflags="-s -w" -o server cmd/server/main.go

# This reduces binary size by ~20-30%
```

### **3. Hot Reload (Development)**
```bash
# Install air for hot reloading
go install github.com/cosmtrek/air@latest

# Create .air.toml configuration
# Run with hot reload
air
```

## 🔄 **Development Workflow**

### **Typical Development Cycle**
```bash
# 1. Make code changes
# 2. Test with go run
go run cmd/server/main.go

# 3. If successful, build for testing
go build -o server cmd/server/main.go

# 4. Test the built version
./server

# 5. Repeat cycle
```

## 📱 **Platform-Specific Notes**

### **macOS**
```bash
# Build for macOS
go build -o server-macos cmd/server/main.go

# Run on macOS
./server-macos
```

### **Linux**
```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o server-linux cmd/server/main.go

# Transfer to Linux server
scp server-linux user@server:/path/to/app/

# Run on Linux
./server-linux
```

### **Windows**
```bash
# Build for Windows
GOOS=windows GOARCH=amd64 go build -o server.exe cmd/server/main.go

# Run on Windows
server.exe
```

## ✅ **Success Criteria**

Your backend is successfully built and running when:

1. **✅ Build completes** without errors
2. **✅ Server starts** on port 8080
3. **✅ Health endpoint** responds correctly
4. **✅ Database connection** established
5. **✅ API endpoints** accessible
6. **✅ No error logs** in console

## 🎯 **Next Steps After Successful Build**

1. **Test API endpoints** with curl or Postman
2. **Start frontend** development server
3. **Test complete flow** - Backend → Frontend
4. **Verify data persistence** in Firestore

---

**Ready to build!** 🚀 Follow this guide to get your backend running successfully.
