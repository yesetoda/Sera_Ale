package main

import (
	"fmt"
	"log"
	"os"

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
// @host localhost:8080
// @BasePath /
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

	// Set up Gin
	r := gin.Default()
	// Swagger docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(gsfiles.Handler))
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public routes
	r.POST("/signup", userHandler.Signup)
	r.POST("/login", userHandler.Login)

	// Auth middleware
	auth := middleware.AuthMiddleware(os.Getenv("JWT_SECRET"))

	// Company routes
	company := r.Group("/company", auth, middleware.RequireRole("company"))
	company.POST("/jobs", jobHandler.CreateJob)
	company.PUT("/jobs/:id", jobHandler.UpdateJob)
	company.DELETE("/jobs/:id", jobHandler.DeleteJob)
	company.GET("/jobs", jobHandler.GetJobsByCompany)
	company.GET("/applications/job", appHandler.GetApplicationsForJob)
	company.PUT("/applications/:id/status", appHandler.UpdateStatus)

	// Applicant routes
	applicant := r.Group("/applicant", auth, middleware.RequireRole("applicant"))
	applicant.GET("/jobs", jobHandler.SearchJobs)
	applicant.GET("/jobs/:id", jobHandler.GetJob)
	applicant.POST("/applications", appHandler.Apply)
	applicant.GET("/applications", appHandler.TrackApplications)

	// Authenticated (any role)
	authd := r.Group("/jobs", auth)
	authd.GET(":id", jobHandler.GetJob)

	fmt.Println("Server running on :8080")
	r.Run(":8080")
}
