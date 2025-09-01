# 🧪 Backend Testing Guide

## 📋 Prerequisites Check

Before testing, ensure you have:

1. **✅ .env file** - Created with correct values
2. **✅ Firebase service account JSON** - Downloaded and placed in `configs/`
3. **✅ Go installed** - `go version` should work
4. **✅ Dependencies installed** - `go mod tidy` completed

## 🚀 Step-by-Step Testing

### **Phase 1: Environment Setup Verification**

#### 1.1 Check Environment Variables
```bash
# Verify .env file exists and has correct values
cat .env

# Expected content:
# FIRESTORE_PROJECT_ID=your-project-id
# GOOGLE_APPLICATION_CREDENTIALS=./configs/firebase-service-account.json
# JWT_SECRET=your-secret-key
# DB_PROVIDER=firestore
```

#### 1.2 Verify Firebase Credentials
```bash
# Check if service account file exists
ls -la configs/firebase-service-account.json

# File should exist and be readable
# Size should be ~1-2KB
```

#### 1.3 Check Go Dependencies
```bash
# Ensure all dependencies are downloaded
go mod tidy

# Should complete without errors
```

### **Phase 2: Backend Startup Testing**

#### 2.1 Start the Backend Server
```bash
# Start the server
go run cmd/server/main.go
```

#### 2.2 Expected Startup Output
```
[INFO] Starting server on port 8080
[INFO] Using Firestore provider
[INFO] Firestore connection established
[INFO] Server started successfully
```

#### 2.3 If Errors Occur
- **"Project ID not found"** → Check `FIRESTORE_PROJECT_ID` in `.env`
- **"Credentials file not found"** → Verify `firebase-service-account.json` path
- **"Permission denied"** → Check Firebase service account permissions
- **"Port already in use"** → Kill process on port 8080

### **Phase 3: API Endpoint Testing**

#### 3.1 Health Check Endpoint
```bash
# Test basic connectivity
curl http://localhost:8080/api/v1/health

# Expected response:
{
  "status": "healthy",
  "timestamp": "2024-01-01T12:00:00Z",
  "service": "admin-portal-backend"
}
```

#### 3.2 Authentication Endpoints
```bash
# Test registration endpoint
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "testpassword123",
    "name": "Test User"
  }'

# Expected: 201 Created with user data and token
```

#### 3.3 Protected Endpoints
```bash
# Get the token from registration response
TOKEN="your-jwt-token-here"

# Test protected company endpoint
curl -X GET http://localhost:8080/api/v1/companies/me \
  -H "Authorization: Bearer $TOKEN"

# Expected: 200 OK with company data or 404 if no company exists
```

### **Phase 4: Database Connectivity Testing**

#### 4.1 Firestore Connection Test
```bash
# Test database operations
curl -X POST http://localhost:8080/api/v1/companies \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Test Company",
    "domain": "testcompany.com",
    "color_theme": "#3B82F6"
  }'

# Expected: 201 Created with company data
```

#### 4.2 Data Retrieval Test
```bash
# Test getting the company we just created
curl -X GET http://localhost:8080/api/v1/companies/me \
  -H "Authorization: Bearer $TOKEN"

# Expected: 200 OK with company data
```

### **Phase 5: Error Handling Testing**

#### 5.1 Invalid Token Test
```bash
# Test with invalid token
curl -X GET http://localhost:8080/api/v1/companies/me \
  -H "Authorization: Bearer invalid-token"

# Expected: 401 Unauthorized
```

#### 5.2 Missing Token Test
```bash
# Test without token
curl -X GET http://localhost:8080/api/v1/companies/me

# Expected: 401 Unauthorized
```

#### 5.3 Invalid Data Test
```bash
# Test with invalid company data
curl -X POST http://localhost:8080/api/v1/companies \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "",
    "domain": "invalid-domain"
  }'

# Expected: 400 Bad Request with validation errors
```

## 🔧 **Testing Tools**

### **Option 1: Command Line (curl)**
```bash
# Install curl if not available
# macOS: brew install curl
# Ubuntu: sudo apt-get install curl

# Basic curl usage
curl -X METHOD http://localhost:8080/endpoint \
  -H "Header: value" \
  -d '{"key": "value"}'
```

### **Option 2: Postman**
1. Download [Postman](https://www.postman.com/)
2. Create a new collection
3. Set base URL: `http://localhost:8080/api/v1`
4. Create requests for each endpoint

### **Option 3: Browser Developer Tools**
1. Open browser console
2. Use `fetch()` API
3. Test endpoints directly

## 📊 **Expected Test Results**

| Test | Expected Result | Status |
|------|----------------|---------|
| Server startup | No errors, port 8080 listening | ⏳ |
| Health endpoint | 200 OK with service info | ⏳ |
| Registration | 201 Created with user data | ⏳ |
| Login | 200 OK with token | ⏳ |
| Protected endpoint | 200 OK with data | ⏳ |
| Invalid token | 401 Unauthorized | ⏳ |
| Database operations | Data saved/retrieved | ⏳ |

## 🐛 **Common Issues & Solutions**

### **Issue 1: "Port 8080 already in use"**
```bash
# Find process using port 8080
lsof -ti:8080

# Kill the process
kill -9 $(lsof -ti:8080)
```

### **Issue 2: "Firestore connection failed"**
- Verify `FIRESTORE_PROJECT_ID` in `.env`
- Check `firebase-service-account.json` exists
- Ensure service account has Firestore permissions

### **Issue 3: "CORS errors"**
- Verify `CORS_ORIGINS` in `.env` includes frontend URL
- Frontend runs on `http://localhost:5173`

### **Issue 4: "JWT errors"**
- Generate new `JWT_SECRET` using:
  ```bash
  openssl rand -base64 32
  ```

## ✅ **Success Criteria**

Your backend is working correctly when:

1. **✅ Server starts without errors**
2. **✅ Health endpoint responds correctly**
3. **✅ Authentication endpoints work**
4. **✅ Protected endpoints require valid tokens**
5. **✅ Database operations succeed**
6. **✅ Error handling works properly**

## 🚀 **Next Steps After Testing**

Once backend testing passes:

1. **Start frontend**: `cd ../frontend && npm run dev`
2. **Test complete flow**: Login → Setup → Dashboard
3. **Verify data persistence** in Firestore console
4. **Test user management** features

## 📚 **Additional Resources**

- [Go HTTP Testing](https://golang.org/pkg/net/http/httptest/)
- [Firebase Admin SDK](https://firebase.google.com/docs/admin/setup)
- [JWT Testing](https://jwt.io/)

---

**Ready to test!** 🧪 Follow this guide step by step to verify your backend connectivity and functionality.
