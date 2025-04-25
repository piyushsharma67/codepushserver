package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/piyushsharma67/codepushserver/database"
)

type UserHandler struct {
	db database.Database
}

func NewUserHandler(db database.Database) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.db.FindUserByID(uint(userID.(float64)))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      user.ID,
		"email":   user.Email,
		"app_id":  user.AppID,
		"token":   user.Token,
	})
}

type CreateAppRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (h *UserHandler) CreateApp(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req CreateAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.db.FindUserByID(uint(userID.(float64)))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update user's app information
	user.AppID = generateRandomString(32)
	user.Token = generateRandomString(64)

	if err := h.db.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create app"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "App created successfully",
		"app_id":  user.AppID,
		"token":   user.Token,
	})
}

func (h *UserHandler) GetApp(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	appID := c.Param("id")
	user, err := h.db.FindUserByAppID(appID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "App not found"})
		return
	}

	if user.ID != uint(userID.(float64)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"app_id":  user.AppID,
		"token":   user.Token,
	})
}

func (h *UserHandler) UpdateApp(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	appID := c.Param("id")
	user, err := h.db.FindUserByAppID(appID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "App not found"})
		return
	}

	if user.ID != uint(userID.(float64)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Generate new token
	user.Token = generateRandomString(64)

	if err := h.db.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update app"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "App updated successfully",
		"token":   user.Token,
	})
}

func (h *UserHandler) DeleteApp(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	appID := c.Param("id")
	user, err := h.db.FindUserByAppID(appID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "App not found"})
		return
	}

	if user.ID != uint(userID.(float64)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Reset app ID and token
	user.AppID = ""
	user.Token = ""

	if err := h.db.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete app"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "App deleted successfully",
	})
} 