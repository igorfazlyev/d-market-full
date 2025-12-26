package main

import (
	"dental-marketplace/backend/internal/auth"
	"dental-marketplace/backend/internal/config"
	"dental-marketplace/backend/internal/database"
	"dental-marketplace/backend/internal/handlers"
	"dental-marketplace/backend/internal/middleware"
	"dental-marketplace/backend/internal/models"
	"dental-marketplace/backend/internal/repository"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Connect to database
	db, err := database.Connect(cfg.Database.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations (includes seeding)
	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessExpiry,
		cfg.JWT.RefreshExpiry,
	)

	// Initialize repositories
	repo := repository.NewRepository(db.DB)
	constantsRepo := repository.NewConstantsRepository(db.DB)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(repo, jwtManager)
	patientHandler := handlers.NewPatientHandler(repo)
	clinicHandler := handlers.NewClinicHandler(repo)
	regulatorHandler := handlers.NewRegulatorHandler(repo)

	// Setup router
	router := setupRouter(authHandler, patientHandler, clinicHandler, regulatorHandler, constantsRepo, jwtManager)

	// Print startup information
	printStartupInfo(cfg)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupRouter configures all API routes
func setupRouter(
	authHandler *handlers.AuthHandler,
	patientHandler *handlers.PatientHandler,
	clinicHandler *handlers.ClinicHandler,
	regulatorHandler *handlers.RegulatorHandler,
	constantsRepo *repository.ConstantsRepository,
	jwtManager *auth.JWTManager,
) *gin.Engine {
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RequestLogger())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "dental-marketplace-api",
		})
	})

	// API v1 routes
	api := router.Group("/api")
	{
		// Public routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Constants endpoint (public - no auth required)
		commonHandler := handlers.NewCommonHandler(constantsRepo)
		api.GET("/constants", commonHandler.GetConstants)

		// Protected routes - require authentication
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(jwtManager))
		{
			// Auth routes (authenticated)
			protected.GET("/auth/me", authHandler.GetMe)

			// Patient routes
			patient := protected.Group("/patient")
			patient.Use(middleware.RequireRole(models.RolePatient))
			{
				// Scans
				patient.GET("/scans", patientHandler.GetScans)

				// Scan details - use separate route groups
				scans := patient.Group("/scans/:id")
				{
					scans.GET("", patientHandler.GetScanByID)
					scans.GET("/plan", patientHandler.GetTreatmentPlan)
				}

				// Plans
				patient.GET("/plans", patientHandler.GetTreatmentPlans)

				plans := patient.Group("/plans/:plan_id")
				{
					plans.GET("/offers", patientHandler.GetOffers)
				}

				// Other routes
				patient.POST("/search-criteria", patientHandler.UpdateSearchCriteria)
				patient.POST("/select-offer", patientHandler.SelectOffer)
				patient.GET("/appointments", patientHandler.GetAppointments)
				patient.POST("/reviews", patientHandler.CreateReview)
				patient.POST("/complaints", patientHandler.CreateComplaint)
			}

			// Clinic routes
			clinic := protected.Group("/clinic")
			clinic.Use(middleware.RequireRole(models.RoleClinic))
			{
				clinic.GET("/dashboard", clinicHandler.GetDashboard)
				clinic.GET("/incoming-plans", clinicHandler.GetIncomingPlans)
				clinic.POST("/offers", clinicHandler.CreateOffer)
				clinic.GET("/leads", clinicHandler.GetLeads)
				clinic.GET("/appointments", clinicHandler.GetAppointments)
				clinic.PUT("/appointments/:id", clinicHandler.UpdateAppointment)
				clinic.GET("/price-list", clinicHandler.GetPriceList)
				clinic.PUT("/price-list", clinicHandler.UpdatePriceList)
				clinic.GET("/analytics", clinicHandler.GetAnalytics)
			}

			// Regulator routes
			regulator := protected.Group("/regulator")
			regulator.Use(middleware.RequireRole(models.RoleRegulator))
			{
				regulator.GET("/dashboard", regulatorHandler.GetDashboard)
				regulator.GET("/statistics", regulatorHandler.GetStatistics)
				regulator.GET("/clinics", regulatorHandler.GetClinics)
				regulator.GET("/clinics/:id", regulatorHandler.GetClinicDetails)
				regulator.GET("/complaints", regulatorHandler.GetComplaints)
				regulator.GET("/disease-analytics", regulatorHandler.GetDiseaseAnalytics)
			}
		}
	}

	return router
}

func printStartupInfo(cfg *config.Config) {
	log.Println("")
	log.Println("====================================================")
	log.Println("🦷 Dental Marketplace API Server")
	log.Println("====================================================")
	log.Println("")
	log.Println("📊 Configuration:")
	log.Printf("   Server Port:    %s", cfg.Server.Port)
	log.Printf("   Environment:    %s", cfg.Server.GinMode)
	log.Printf("   Database:       %s:%s", cfg.Database.Host, cfg.Database.Port)
	log.Println("")
	log.Println("🔐 Demo Credentials:")
	log.Println("   Patient:        username: patient   | password: password")
	log.Println("   Clinic 1:       username: clinic1   | password: password")
	log.Println("   Clinic 2:       username: clinic2   | password: password")
	log.Println("   Regulator:      username: regulator | password: password")
	log.Println("")
	log.Println("📡 API Endpoints:")
	log.Printf("   Health Check:   http://localhost:%s/health", cfg.Server.Port)
	log.Printf("   Constants:      http://localhost:%s/api/constants", cfg.Server.Port)
	log.Printf("   Login:          http://localhost:%s/api/auth/login", cfg.Server.Port)
	log.Printf("   API Docs:       http://localhost:%s/api", cfg.Server.Port)
	log.Println("")
	log.Println("✨ Features:")
	log.Println("   ✓ JWT Authentication")
	log.Println("   ✓ Role-based Access Control")
	log.Println("   ✓ Patient Management")
	log.Println("   ✓ Clinic Operations")
	log.Println("   ✓ Regulator Dashboard")
	log.Println("   ✓ Treatment Plans & Offers")
	log.Println("   ✓ Analytics & Statistics")
	log.Println("   ✓ Database-driven Constants")
	log.Println("")
	log.Println("====================================================")
	log.Printf("🚀 Server starting on http://localhost:%s", cfg.Server.Port)
	log.Println("====================================================")
	log.Println("")
}
