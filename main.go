package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	gsfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/yesetoda/Sera_Ale/docs" // Swagger docs import
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/yesetoda/Sera_Ale/internal/app"
	"github.com/yesetoda/Sera_Ale/internal/handler"
	"github.com/yesetoda/Sera_Ale/internal/middleware"
	"github.com/yesetoda/Sera_Ale/internal/repository"
	"github.com/yesetoda/Sera_Ale/internal/service"
)

// @title Sera Ale Job Board API
// @version 1.0
// @description RESTful API for job board with applicants and companies
// @host https://sera-ale-4.onrender.com
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Provide your JWT token in the format: Bearer {token}
func main() {
	// Load env vars
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env vars")
	}
	// Connect to DB
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}
	// Ping DB
	sqlDB, _ := db.DB()
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("failed to ping database: ", err)
	}

	// Dependency injection
	userRepo := repository.NewUserRepository(db)
	jobRepo := repository.NewJobRepository(db)
	appRepo := repository.NewApplicationRepository(db)
	jwtSvc := service.NewJWTService()
	pwdSvc := service.NewPasswordService()
	cloudSvc, err := service.NewCloudinaryService()
	if err != nil {
		log.Fatal("failed to init cloudinary: ", err)
	}
	userApp := app.NewUserApp(userRepo, jwtSvc, pwdSvc)
	jobApp := app.NewJobApp(jobRepo)
	appApp := app.NewApplicationApp(appRepo, jobRepo, cloudSvc)

	userHandler := handler.NewUserHandler(userApp)
	jobHandler := handler.NewJobHandler(jobApp)
	appHandler := handler.NewApplicationHandler(appApp)
	authHandler := handler.NewAuthHandler(userApp)

	// Set up Gin
	r := gin.Default()

	// CORS middleware (allow all for dev)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Swagger docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(gsfiles.Handler))

	// Default root route for API help
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":   "Welcome to the Sera Ale Job Board API! See /swagger/index.html for documentation.",
			"docs":      "/swagger/index.html",
			"endpoints": []string{"/signup", "/login", "/user/me", "/company/jobs", "/applicant/jobs", "/applicant/applications", "/company/applications/job"},
		})
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public routes
	r.POST("/signup", authHandler.Signup)
	r.POST("/login", authHandler.Login)
	r.POST("/auth/signup", authHandler.Signup)
	r.POST("/auth/login", authHandler.Login)

	// Auth middleware
	auth := middleware.AuthMiddleware(os.Getenv("JWT_SECRET"))

	// Company routes
	r.GET("/jobs", jobHandler.GetJobsByCompany)
	// Requires Bearer token in Authorization header.
	company := r.Group("/company", auth, middleware.RequireRole("company"))
	company.POST("/jobs", jobHandler.CreateJob)
	company.PUT("/jobs/:id", jobHandler.UpdateJob)
	company.DELETE("/jobs/:id", jobHandler.DeleteJob)
	company.GET("/applications/job", appHandler.GetApplicationsForJob)
	company.PUT("/applications/:id/status", appHandler.UpdateStatus)

	// Applicant routes
	// Requires Bearer token in Authorization header.
	applicant := r.Group("/applicant", auth, middleware.RequireRole("applicant"))
	applicant.GET("/jobs", jobHandler.SearchJobs)
	applicant.GET("/jobs/:id", jobHandler.GetJob)
	applicant.POST("/applications", appHandler.Apply)
	applicant.GET("/applications", appHandler.TrackApplications)

	// Public job details route (matches @Router /jobs/{id} [get])
	r.GET("/jobs/:id", jobHandler.GetJob)

	// User profile route
	// Requires Bearer token in Authorization header.
	r.GET("/user/me", auth, userHandler.GetCurrentUser)

	// 404 Not Found handler
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"success": false, "message": "Route not found. See /swagger/index.html for API documentation."})
	})

	// 405 Method Not Allowed handler
	r.NoMethod(func(c *gin.Context) {
		c.JSON(405, gin.H{"success": false, "message": "Method not allowed on this route."})
	})

	fmt.Println("Server running on :8080")
	r.Run(":8080")
}
