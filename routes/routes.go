package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/piyushsharma67/codepushserver/config"
	"github.com/piyushsharma67/codepushserver/database"
	"github.com/piyushsharma67/codepushserver/handlers"
	"github.com/piyushsharma67/codepushserver/middleware"
)

func SetupRoutes(router *gin.Engine) {
	// Initialize handlers
	db, err := database.NewDatabase(config.LoadConfig())
	if err != nil {
		panic("Failed to initialize database")
	}

	authHandler := handlers.NewAuthHandler(db)
	userHandler := handlers.NewUserHandler(db)

	// Public routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Protected routes
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/profile", userHandler.GetProfile)
		api.POST("/app", userHandler.CreateApp)
		api.GET("/app/:id", userHandler.GetApp)
		api.PUT("/app/:id", userHandler.UpdateApp)
		api.DELETE("/app/:id", userHandler.DeleteApp)
	}
} 