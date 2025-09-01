package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mohilkhare21/pab-smb-poc-golang/backend/internal/auth"
	"github.com/mohilkhare21/pab-smb-poc-golang/backend/internal/database"
	"github.com/mohilkhare21/pab-smb-poc-golang/backend/internal/handlers"
	"github.com/mohilkhare21/pab-smb-poc-golang/backend/internal/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database provider
	dbConfig := database.DatabaseConfig{
		Provider:    getEnv("DB_PROVIDER", "firestore"),
		ProjectID:   getEnv("FIRESTORE_PROJECT_ID", ""),
		Host:        getEnv("DB_HOST", ""),
		Port:        getEnvAsInt("DB_PORT", 0),
		Username:    getEnv("DB_USERNAME", ""),
		Password:    getEnv("DB_PASSWORD", ""),
		Database:    getEnv("DB_NAME", ""),
		SSLMode:     getEnv("DB_SSL_MODE", ""),
		Credentials: getEnv("GOOGLE_APPLICATION_CREDENTIALS", ""),
	}

	dbFactory := &database.DefaultDatabaseFactory{}
	dbProvider, err := dbFactory.CreateProvider(dbConfig)
	if err != nil {
		log.Fatalf("Failed to create database provider: %v", err)
	}
	defer dbProvider.Close()

	// Test database connection
	if err := dbProvider.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connection established")

	// Initialize auth provider
	authConfig := auth.AuthConfig{
		Provider:     getEnv("AUTH_PROVIDER", "auth0"),
		Domain:       getEnv("AUTH0_DOMAIN", ""),
		ClientID:     getEnv("AUTH0_CLIENT_ID", ""),
		ClientSecret: getEnv("AUTH0_CLIENT_SECRET", ""),
		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key"),
		RedirectURL:  getEnv("AUTH0_REDIRECT_URL", ""),
	}

	authFactory := &auth.DefaultAuthFactory{}
	authProvider, err := authFactory.CreateProvider(authConfig)
	if err != nil {
		log.Fatalf("Failed to create auth provider: %v", err)
	}

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authProvider, dbProvider)
	companyHandler := handlers.NewCompanyHandler(dbProvider)
	userHandler := handlers.NewUserHandler(dbProvider, authProvider)
	invitationHandler := handlers.NewInvitationHandler(dbProvider, authProvider)
	shortcutHandler := handlers.NewBrowserShortcutHandler(dbProvider)
	setupHandler := handlers.NewSetupHandler(dbProvider)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authProvider)

	// Set up Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Add CORS middleware
	router.Use(middleware.CORS())

	// Add request logging middleware
	router.Use(middleware.RequestLogger())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"service":   "multi-tenant-admin-portal",
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		public := api.Group("")
		{
			public.POST("/auth/login", authHandler.Login)
			public.POST("/auth/register", authHandler.Register)
			public.POST("/auth/reset-password", authHandler.ResetPassword)
			public.GET("/auth/verify", authHandler.AuthenticateWithToken)
		}

		// Protected routes (authentication required)
		protected := api.Group("")
		protected.Use(authMiddleware.Authenticate())
		{
			// Auth routes
			protected.POST("/auth/logout", authHandler.Logout)
			protected.POST("/auth/refresh", authHandler.RefreshToken)
			protected.POST("/auth/change-password", authHandler.ChangePassword)

			// Company routes
			protected.POST("/companies", companyHandler.CreateCompany)
			protected.GET("/companies/me", companyHandler.GetCompany)
			protected.PUT("/companies/me", companyHandler.UpdateCompany)
			protected.DELETE("/companies/me", companyHandler.DeleteCompany)
			protected.GET("/companies/stats", companyHandler.GetCompanyStats)

			// User routes
			protected.GET("/users", userHandler.GetUsers)
			protected.GET("/users/:id", userHandler.GetUser)
			protected.PUT("/users/:id", userHandler.UpdateUser)
			protected.DELETE("/users/:id", userHandler.DeleteUser)

			// Invitation routes
			protected.POST("/invitations", invitationHandler.CreateInvitation)
			protected.GET("/invitations", invitationHandler.GetInvitations)
			protected.DELETE("/invitations/:id", invitationHandler.DeleteInvitation)
			protected.POST("/invitations/:token/accept", invitationHandler.AcceptInvitation)

			// Browser shortcut routes
			protected.GET("/shortcuts", shortcutHandler.GetShortcuts)
			protected.POST("/shortcuts", shortcutHandler.CreateShortcut)
			protected.PUT("/shortcuts/:id", shortcutHandler.UpdateShortcut)
			protected.DELETE("/shortcuts/:id", shortcutHandler.DeleteShortcut)

			// Setup and configuration routes
			protected.GET("/setup/progress", setupHandler.GetSetupProgress)
			protected.PUT("/setup/step", setupHandler.UpdateSetupStep)
			protected.GET("/setup/stats", setupHandler.GetCompanyStats)
			protected.PUT("/setup/config", setupHandler.UpdateConfigurationStatus)
			protected.POST("/setup/generate-shortcuts", setupHandler.GenerateShortcuts)
			protected.POST("/setup/nudge-users", setupHandler.NudgeUsers)
			protected.GET("/setup/download-info", setupHandler.GetDownloadInfo)
		}

		// Admin routes (admin role required)
		admin := api.Group("/admin")
		admin.Use(authMiddleware.Authenticate())
		admin.Use(authMiddleware.RequireRole("admin"))
		{
			admin.GET("/companies", companyHandler.ListCompanies)
		}
	}

	// Get server configuration
	port := getEnv("PORT", "8080")
	serverAddr := fmt.Sprintf(":%s", port)

	// Create HTTP server
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	log.Printf("Health check available at http://localhost:%s/health", port)
	log.Printf("API documentation available at http://localhost:%s/api/v1", port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := fmt.Sscanf(value, "%d", &defaultValue); err == nil && intValue == 1 {
			return defaultValue
		}
	}
	return defaultValue
}

