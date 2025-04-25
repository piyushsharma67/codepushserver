package main

import (
	"log"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/piyushsharma67/codepushserver/config"
	"github.com/piyushsharma67/codepushserver/database"
	"github.com/piyushsharma67/codepushserver/handlers"
	v2 "github.com/piyushsharma67/codepushserver/handlers/v2"
	store "github.com/piyushsharma67/codepushserver/store/v2"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize store
	store := store.NewStore(db)

	// Initialize router
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// V1 routes
	router.POST("/auth/register", handlers.NewAuthHandler(db).Register)
	router.POST("/auth/login", handlers.NewAuthHandler(db).Login)

	// V2 routes
	v2Router := v2.SetupRouter(store)
	router.Any("/v2/*path", func(c *gin.Context) {
		v2Router.HandleContext(c)
	})

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// TODO: Implement JWT validation
		// For now, just set a dummy user ID
		c.Set("user_id", uint(1))
		c.Next()
	}
} 