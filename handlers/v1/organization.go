package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/piyushsharma67/codepushserver/database"
	v1 "github.com/piyushsharma67/codepushserver/services/v1"
	"github.com/piyushsharma67/codepushserver/utils"
)

type OrganizationHandler struct {
	orgService *v1.OrganizationService
}

func NewOrganizationHandler(db database.Database) *OrganizationHandler {
	return &OrganizationHandler{
		orgService: v1.NewOrganizationService(db),
	}
}

type CreateOrganizationRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org, err := h.orgService.CreateOrganization(userID, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organization"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Organization created successfully",
		"organization": gin.H{
			"id":           org.ID,
			"name":         org.Name,
			"description":  org.Description,
			"public_token": org.PublicToken,
			"created_at":   org.CreatedAt,
		},
	})
}

func (h *OrganizationHandler) GetOrganization(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	org, err := h.orgService.GetOrganization(userID, orgID)
	if err != nil {
		if err == utils.ErrAccessDenied {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           org.ID,
		"name":         org.Name,
		"description":  org.Description,
		"public_token": org.PublicToken,
		"created_at":   org.CreatedAt,
	})
}

func (h *OrganizationHandler) GetUserOrganizations(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	orgs, err := h.orgService.GetUserOrganizations(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch organizations"})
		return
	}

	var responseOrgs []gin.H
	for _, org := range orgs {
		responseOrgs = append(responseOrgs, gin.H{
			"id":           org.ID,
			"name":         org.Name,
			"description":  org.Description,
			"public_token": org.PublicToken,
			"created_at":   org.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"organizations": responseOrgs,
	})
}

type InviteUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required"`
}

func (h *OrganizationHandler) InviteUser(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	var req InviteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.orgService.InviteUser(userID, orgID, req.Email, req.Role)
	if err != nil {
		if err == utils.ErrAccessDenied {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to invite user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User invited successfully",
	})
}

func (h *OrganizationHandler) GetPendingInvites(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	invites, err := h.orgService.GetPendingInvites(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pending invites"})
		return
	}

	var responseInvites []gin.H
	for _, invite := range invites {
		responseInvites = append(responseInvites, gin.H{
			"id":             invite.ID,
			"organization_id": invite.OrganizationID,
			"email":          invite.Email,
			"role":           invite.Role,
			"created_at":     invite.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"invites": responseInvites,
	})
}

func (h *OrganizationHandler) DeleteOrganization(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	err = h.orgService.DeleteOrganization(userID, orgID)
	if err != nil {
		if err == utils.ErrAccessDenied {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete organization"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Organization deleted successfully",
	})
}

type TransferAdminRequest struct {
	NewAdminID uint `json:"new_admin_id" binding:"required"`
}

func (h *OrganizationHandler) TransferAdmin(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	var req TransferAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.orgService.TransferAdmin(userID, orgID, req.NewAdminID)
	if err != nil {
		if err == utils.ErrAccessDenied {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to transfer admin role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Admin role transferred successfully",
	})
} 