package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"backend/config"
	"backend/handlers"
	"backend/logger"
	"backend/middleware"
	"backend/repository"
	"backend/routes"
	"backend/seed"
	"backend/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// ── Environment ─────────────────────────────────────────────
	if err := godotenv.Load(); err != nil {
		// not fatal – we fall back to system env vars
	}

	requiredEnvVars := []string{"MONGO_URI", "DB_NAME", "JWT_SECRET"}
	for _, v := range requiredEnvVars {
		if os.Getenv(v) == "" {
			panic("FATAL: required environment variable " + v + " is not set")
		}
	}
	if len(os.Getenv("JWT_SECRET")) < 32 {
		panic("FATAL: JWT_SECRET must be at least 32 characters long")
	}

	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)

	// ── Structured Logger ───────────────────────────────────────
	logger.Init(ginMode)

	// ── Database ────────────────────────────────────────────────
	config.ConnectDB()
	seed.SeedData()

	// ── Dependency Injection ────────────────────────────────────
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

	// ── Router & Middleware ─────────────────────────────────────
	r := gin.New()

	// 1. Panic recovery
	r.Use(gin.Recovery())

	// 2. Request ID
	r.Use(middleware.RequestID())

	// 3. Rate limiter
	r.Use(middleware.RateLimiter())

	// 4. Structured request logging (replaces gin.Logger)
	r.Use(middleware.RequestLogger())

	// 5. Request timeout
	r.Use(middleware.Timeout())

	// 6. GZIP compression
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// 7. CORS
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
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// 8. Security headers
	r.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	})

	// ── Routes ──────────────────────────────────────────────────
	routes.SetupRoutes(r, h)

	// Liveness probe – always OK if process is running
	r.GET("/health/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "alive"})
	})

	// Readiness probe – checks database connectivity
	r.GET("/health/ready", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()
		if err := config.PingDB(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	// Legacy health endpoint (backward compatible)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Travel Planner API is running"})
	})

	// ── Server ──────────────────────────────────────────────────
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	srv := &http.Server{
		Addr:           ":" + port,
		Handler:        r,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// ── Graceful Shutdown ───────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Log.Info().Str("port", port).Str("mode", ginMode).Msg("Server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	<-quit
	logger.Log.Info().Msg("Shutdown signal received, draining connections…")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Error().Err(err).Msg("Forced server shutdown")
	}

	config.DisconnectDB()
	logger.Log.Info().Msg("Server exited gracefully")
}
