package v2

import (
	"github.com/gin-gonic/gin"
	"github.com/piyushsharma67/codepushserver/middleware"
	service "github.com/piyushsharma67/codepushserver/services/v2"
	store "github.com/piyushsharma67/codepushserver/store/v2"
)

func SetupRouter(store store.Store) *gin.Engine {
	router := gin.Default()

	// Initialize services
	authService := service.NewAuthService(store)
	userService := service.NewUserService(store)

	// Initialize handlers
	authHandler := NewAuthHandler(authService)
	userHandler := NewUserController(userService)
	healthHandler := NewHealthHandler()

	// Public routes
	router.POST("/v2/auth/register", authHandler.Register)
	router.POST("/v2/auth/login", authHandler.Login)

	// Health check route
	router.GET("/v2/health", healthHandler.Check)

	// Protected routes
	v2 := router.Group("/v2")
	v2.Use(middleware.AuthMiddleware())
	{
		// User routes
		v2.GET("/user/profile", userHandler.GetProfile)
		v2.PUT("/user/profile", userHandler.UpdateProfile)

		// App routes
		v2.POST("/apps", userHandler.CreateApp)
		v2.GET("/apps", userHandler.ListApps)
		v2.GET("/apps/:app_id", userHandler.GetApp)
		v2.PUT("/apps/:app_id", userHandler.UpdateApp)
		v2.DELETE("/apps/:app_id", userHandler.DeleteApp)
	}

	return router
} 