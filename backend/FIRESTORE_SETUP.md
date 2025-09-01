# 🔥 Firestore Setup Guide

## 📋 Prerequisites

1. **Google Cloud Project** - You mentioned you already created one
2. **Firebase Console Access** - https://console.firebase.google.com/
3. **Service Account Key** - JSON file with credentials

## 🚀 Step-by-Step Setup

### Step 1: Get Your Firebase Project ID

1. Go to [Firebase Console](https://console.firebase.google.com/)
2. Select your project
3. Click on the gear icon (⚙️) next to "Project Overview"
4. Copy the **Project ID** (e.g., `my-awesome-project-123`)

### Step 2: Create Service Account Key

1. In Firebase Console, go to **Project Settings**
2. Click on **Service Accounts** tab
3. Click **Generate new private key**
4. Download the JSON file
5. **Save it as** `firebase-service-account.json` in `backend/configs/` directory

### Step 3: Set Up Environment Variables

1. **Rename** `env.firestore` to `.env` in the backend directory:
   ```bash
   cd backend
   mv env.firestore .env
   ```

2. **Edit** the `.env` file with your actual values:
   ```bash
   # Replace these values with your actual Firebase project details
   FIRESTORE_PROJECT_ID=your-actual-project-id-here
   GOOGLE_APPLICATION_CREDENTIALS=./configs/firebase-service-account.json
   JWT_SECRET=generate-a-random-secret-key-here
   ```

### Step 4: Generate JWT Secret

Generate a secure JWT secret:
```bash
# Option 1: Using openssl
openssl rand -base64 32

# Option 2: Using Python
python3 -c "import secrets; print(secrets.token_urlsafe(32))"

# Option 3: Online generator (for development only)
# https://generate-secret.vercel.app/32
```

### Step 5: Verify File Structure

Your backend directory should look like this:
```
backend/
├── .env                                    # Your environment variables
├── configs/
│   ├── firebase-service-account.json      # Firebase credentials
│   └── env.example                        # Example file
├── go.mod
└── ...
```

## 🔑 Required Environment Variables

| Variable | Description | Example Value |
|----------|-------------|---------------|
| `FIRESTORE_PROJECT_ID` | Your Firebase project ID | `my-project-123` |
| `GOOGLE_APPLICATION_CREDENTIALS` | Path to service account JSON | `./configs/firebase-service-account.json` |
| `JWT_SECRET` | Secret key for JWT tokens | `abc123...` (32+ characters) |
| `DB_PROVIDER` | Database provider | `firestore` |
| `PORT` | Server port | `8080` |
| `CORS_ORIGINS` | Allowed frontend origins | `http://localhost:5173` |

## 📁 Service Account File Location

**Important**: Place your `firebase-service-account.json` in:
```
backend/configs/firebase-service-account.json
```

**DO NOT** commit this file to git! It contains sensitive credentials.

## 🚫 Security Notes

1. **Never commit** `.env` or service account files to git
2. **Add to .gitignore**:
   ```
   .env
   configs/*.json
   ```
3. **Keep credentials secure** - don't share or expose them

## 🧪 Testing the Setup

### 1. Start the Backend
```bash
cd backend
go run cmd/server/main.go
```

### 2. Check Logs
You should see:
```
[INFO] Starting server on port 8080
[INFO] Using Firestore provider
[INFO] Firestore connection established
```

### 3. Test Health Endpoint
```bash
curl http://localhost:8080/api/v1/health
```

Expected response:
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T12:00:00Z",
  "service": "admin-portal-backend"
}
```

## 🐛 Troubleshooting

### Common Issues:

#### 1. **"Project ID not found"**
- Verify `FIRESTORE_PROJECT_ID` in `.env`
- Check the project ID in Firebase Console

#### 2. **"Credentials file not found"**
- Verify `GOOGLE_APPLICATION_CREDENTIALS` path
- Ensure `firebase-service-account.json` exists in `configs/`

#### 3. **"Permission denied"**
- Check if service account has Firestore permissions
- Verify the service account JSON is valid

#### 4. **"CORS error"**
- Verify `CORS_ORIGINS` includes your frontend URL
- Frontend runs on `http://localhost:5173` (Vite default)

## 🔄 Next Steps

Once Firestore is working:

1. **Test backend endpoints** with Postman or curl
2. **Start frontend** with `npm run dev`
3. **Test complete flow** - login → setup → dashboard
4. **Verify data persistence** in Firestore console

## 📚 Additional Resources

- [Firebase Console](https://console.firebase.google.com/)
- [Firestore Documentation](https://firebase.google.com/docs/firestore)
- [Service Account Setup](https://firebase.google.com/docs/admin/setup#initialize-sdk)

## ✅ Checklist

- [ ] Downloaded service account JSON
- [ ] Created `.env` file with correct values
- [ ] Placed service account file in `configs/`
- [ ] Generated secure JWT secret
- [ ] Backend starts without errors
- [ ] Health endpoint responds correctly
- [ ] Frontend can connect to backend

---

**Ready to test!** 🚀 Once you complete these steps, your backend will be fully connected to Firestore and ready to work with the frontend.
