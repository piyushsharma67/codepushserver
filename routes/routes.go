package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/piyushsharma67/codepushserver/database"
	v1 "github.com/piyushsharma67/codepushserver/handlers/v1"
	"github.com/piyushsharma67/codepushserver/middleware"
)

func SetupRoutes(router *gin.Engine, db database.Database) {
	// Initialize handlers
	authHandler := v1.NewAuthHandler(db)
	userHandler := v1.NewUserHandler(db)
	orgHandler := v1.NewOrganizationHandler(db)

	// API v1 routes
	v1Group := router.Group("/api/v1")
	{
		// Auth routes (public)
		v1Group.POST("/auth/register", authHandler.Register)
		v1Group.POST("/auth/login", authHandler.Login)

		// Protected routes
		protected := v1Group.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// User routes
			protected.GET("/user/profile", userHandler.GetProfile)
			protected.PUT("/user/profile", userHandler.UpdateProfile)
			protected.GET("/user/apps", userHandler.GetAllApps)
			protected.POST("/user/apps", userHandler.CreateApp)
			protected.GET("/user/apps/:id", userHandler.GetApp)
			protected.PUT("/user/apps/:id", userHandler.UpdateApp)
			protected.DELETE("/user/apps/:id", userHandler.DeleteApp)

			// Organization routes
			protected.POST("/organizations", orgHandler.CreateOrganization)
			protected.POST("/organizations/invite", orgHandler.InviteUser)
			// protected.POST("/organizations/accept-invite", orgHandler.AcceptInvite)
			protected.GET("/organizations", orgHandler.GetUserOrganizations)
			protected.GET("/organizations/pending-invites", orgHandler.GetPendingInvites)
			protected.DELETE("/organizations/:id", orgHandler.DeleteOrganization)
			protected.POST("/organizations/transfer-admin", orgHandler.TransferAdmin)
		}
	}
} 