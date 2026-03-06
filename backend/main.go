package main

import (
	"log"
	"os"
	"strings"
	"time"

	"backend/config"
	"backend/handlers"
	"backend/repository"
	"backend/routes"
	"backend/seed"
	"backend/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	requiredEnvVars := []string{"MONGO_URI", "DB_NAME", "JWT_SECRET"}
	for _, v := range requiredEnvVars {
		if os.Getenv(v) == "" {
			log.Fatalf("FATAL: required environment variable %s is not set", v)
		}
	}

	if len(os.Getenv("JWT_SECRET")) < 32 {
		log.Fatal("FATAL: JWT_SECRET must be at least 32 characters long")
	}

	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)

	config.ConnectDB()
	seed.SeedData()

	// Dependency Injection
	userRepo := repository.NewUserRepository()
	destRepo := repository.NewDestinationRepository()
	reviewRepo := repository.NewReviewRepository()
	itinRepo := repository.NewItineraryRepository()

	authService := service.NewAuthService(userRepo)
	destService := service.NewDestinationService(destRepo, reviewRepo)
	reviewService := service.NewReviewService(reviewRepo, destRepo, userRepo)
	itinService := service.NewItineraryService(itinRepo, destRepo)

	h := &routes.Handlers{
		Auth:        handlers.NewAuthHandler(authService),
		Destination: handlers.NewDestinationHandler(destService),
		Review:      handlers.NewReviewHandler(reviewService, userRepo),
		Itinerary:   handlers.NewItineraryHandler(itinService),
	}

	r := gin.New()
	r.Use(gin.Recovery())

	if ginMode != gin.ReleaseMode {
		r.Use(gin.Logger())
	}

	allowedOrigins := []string{"http://localhost:3000"}
	if origins := os.Getenv("ALLOWED_ORIGINS"); origins != "" {
		allowedOrigins = strings.Split(origins, ",")
		for i, o := range allowedOrigins {
			allowedOrigins[i] = strings.TrimSpace(o)
		}
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	})

	routes.SetupRoutes(r, h)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "Travel Planner API is running"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("Server starting on port %s (mode: %s)", port, ginMode)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
