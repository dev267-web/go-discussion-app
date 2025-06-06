// main.go 
package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"go-discussion-app/config"
	"go-discussion-app/internal/auth"
	"go-discussion-app/internal/comment"
	"go-discussion-app/internal/discussion"
	"go-discussion-app/internal/health"
	"go-discussion-app/internal/middleware"
	"go-discussion-app/internal/subscription"
	"go-discussion-app/internal/tag"
	"go-discussion-app/internal/user"
	"go-discussion-app/db"
)

func main() {
	// Load environment variables and config
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize DB connection
	dbConn, err := db.InitPostgres(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer dbConn.Close()

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// CORS middleware (allow all for now; restrict in prod)
	router.Use(cors.Default())

	// Global middlewares (e.g., logging)
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Public routes
	auth.RegisterRoutes(router, cfg.JWTSecret)
	health.RegisterRoutes(router)

	// Protected routes group (JWT middleware)
	protected := router.Group("/")
	protected.Use(middleware.JWTAuth(cfg.JWTSecret))

	user.RegisterRoutes(protected, dbConn)
	discussion.RegisterRoutes(protected, dbConn)
	comment.RegisterRoutes(protected, dbConn)
	subscription.RegisterRoutes(protected, dbConn)
	tag.RegisterRoutes(protected, dbConn)

	// Start server
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
