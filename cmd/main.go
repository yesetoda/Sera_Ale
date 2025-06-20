package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	gsfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/swaggo/gin-swagger/example/docs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	// Set up Gin
	r := gin.Default()
	// Swagger docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(gsfiles.Handler))
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	fmt.Println("Server running on :8080")
	r.Run(":8080")
}
