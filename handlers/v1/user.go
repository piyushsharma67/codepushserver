package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/piyushsharma67/codepushserver/database"
	v1 "github.com/piyushsharma67/codepushserver/services/v1"
	"github.com/piyushsharma67/codepushserver/utils"
)

type UserHandler struct {
	userService *v1.UserService
}

func NewUserHandler(db database.Database) *UserHandler {
	return &UserHandler{
		userService: v1.NewUserService(db),
	}
}

type UpdateProfileRequest struct {
	Username    string `json:"username"`
	CompanyName string `json:"company_name"`
	PhoneNumber string `json:"phone_number"`
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.userService.GetUserProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           user.ID,
		"username":     user.Username,
		"email":        user.Email,
		"company_name": user.CompanyName,
		"phone_number": user.PhoneNumber,
		"created_at":   user.CreatedAt,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.UpdateUserProfile(userID, req.Username, req.CompanyName, req.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           user.ID,
		"username":     user.Username,
		"email":        user.Email,
		"company_name": user.CompanyName,
		"phone_number": user.PhoneNumber,
		"created_at":   user.CreatedAt,
	})
}

type CreateAppRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (h *UserHandler) CreateApp(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req CreateAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	app, err := h.userService.CreateApp(userID, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create app"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "App created successfully",
		"app": gin.H{
			"id":          app.ID,
			"name":        app.Name,
			"description": app.Description,
			"token":       app.Token,
			"created_at":  app.CreatedAt,
		},
	})
}

func (h *UserHandler) GetApp(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	appID := c.Param("id")
	app, err := h.userService.GetApp(userID, appID)
	if err != nil {
		if err == utils.ErrAccessDenied {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "App not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          app.ID,
		"name":        app.Name,
		"description": app.Description,
		"token":       app.Token,
		"created_at":  app.CreatedAt,
	})
}

func (h *UserHandler) UpdateApp(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	appID := c.Param("id")
	app, err := h.userService.UpdateApp(userID, appID)
	if err != nil {
		if err == utils.ErrAccessDenied {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update app"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "App updated successfully",
		"app": gin.H{
			"id":          app.ID,
			"name":        app.Name,
			"description": app.Description,
			"token":       app.Token,
			"created_at":  app.CreatedAt,
		},
	})
}

func (h *UserHandler) DeleteApp(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	appID := c.Param("id")
	err := h.userService.DeleteApp(userID, appID)
	if err != nil {
		if err == utils.ErrAccessDenied {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete app"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "App deleted successfully",
	})
}

func (h *UserHandler) GetAllApps(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	apps, err := h.userService.GetAllApps(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch apps"})
		return
	}

	// Convert apps to response format
	var responseApps []gin.H
	for _, app := range apps {
		responseApps = append(responseApps, gin.H{
			"id":          app.ID,
			"name":        app.Name,
			"description": app.Description,
			"token":       app.Token,
			"created_at":  app.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"apps": responseApps,
	})
} 