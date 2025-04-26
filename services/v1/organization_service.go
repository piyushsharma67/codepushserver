package v1

import (
	"github.com/google/uuid"
	"github.com/piyushsharma67/codepushserver/database"
	"github.com/piyushsharma67/codepushserver/models"
	"github.com/piyushsharma67/codepushserver/utils"
)

type OrganizationService struct {
	db database.Database
}

func NewOrganizationService(db database.Database) *OrganizationService {
	return &OrganizationService{db: db}
}

type CreateOrganizationRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type InviteUserRequest struct {
	OrganizationID uint   `json:"organization_id" binding:"required"`
	InviteeEmail   string `json:"invitee_email" binding:"required,email"`
}

type AcceptInviteRequest struct {
	InviteID uint `json:"invite_id" binding:"required"`
}

type TransferAdminRequest struct {
	OrganizationID uint `json:"organization_id" binding:"required"`
	NewAdminID     uint `json:"new_admin_id" binding:"required"`
}

func (s *OrganizationService) CreateOrganization(userID uint, name, description string) (*models.Organization, error) {
	user, err := s.db.FindUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Generate public and private tokens
	publicToken := utils.GenerateRandomString(32)  // Shorter for public use
	privateToken := utils.GenerateRandomString(64) // Longer for secure operations

	org := &models.Organization{
		ID:           uuid.New(),
		Name:         name,
		Description:  description,
		PublicToken:  publicToken,
		PrivateToken: privateToken,
		CreatedBy:    user.ID,
	}

	if err := s.db.CreateOrganization(org); err != nil {
		return nil, err
	}

	// Create organization membership for the creator as admin
	membership := &models.OrganizationMember{
		OrganizationID: org.ID,
		UserID:         user.ID,
		Role:           "admin",
	}

	if err := s.db.CreateOrganizationMember(membership); err != nil {
		return nil, err
	}

	return org, nil
}

func (s *OrganizationService) GetOrganization(userID uint, orgID uuid.UUID) (*models.Organization, error) {
	org, err := s.db.FindOrganizationByID(orgID)
	if err != nil {
		return nil, err
	}

	// Check if user is a member
	_, err = s.db.FindOrganizationMember(orgID, userID)
	if err != nil {
		return nil, utils.ErrAccessDenied
	}

	// Create a copy of the organization without the private token
	orgResponse := *org
	orgResponse.PrivateToken = "" // Don't expose private token

	return &orgResponse, nil
}

func (s *OrganizationService) GetUserOrganizations(userID uint) ([]*models.Organization, error) {
	orgs, err := s.db.FindOrganizationsByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Remove private tokens from all organizations
	for _, org := range orgs {
		org.PrivateToken = ""
	}

	return orgs, nil
}

func (s *OrganizationService) InviteUser(userID uint, orgID uuid.UUID, email string, role string) error {
	// Check if inviter is a member and has admin rights
	member, err := s.db.FindOrganizationMember(orgID, userID)
	if err != nil {
		return utils.ErrAccessDenied
	}

	if member.Role != "admin" {
		return utils.ErrAccessDenied
	}

	// Create invitation
	invitation := &models.OrganizationInvitation{
		OrganizationID: orgID,
		Email:          email,
		Role:           role,
		Status:         "pending",
	}

	return s.db.CreateOrganizationInvitation(invitation)
}

func (s *OrganizationService) GetPendingInvites(userID uint) ([]*models.OrganizationInvitation, error) {
	user, err := s.db.FindUserByID(userID)
	if err != nil {
		return nil, err
	}

	return s.db.FindPendingInvitationsByEmail(user.Email)
}

func (s *OrganizationService) DeleteOrganization(userID uint, orgID uuid.UUID) error {
	// Check if user is admin
	member, err := s.db.FindOrganizationMember(orgID, userID)
	if err != nil {
		return utils.ErrAccessDenied
	}

	if member.Role != "admin" {
		return utils.ErrAccessDenied
	}

	return s.db.DeleteOrganization(orgID)
}

func (s *OrganizationService) TransferAdmin(userID uint, orgID uuid.UUID, newAdminID uint) error {
	// Check if current user is admin
	member, err := s.db.FindOrganizationMember(orgID, userID)
	if err != nil {
		return utils.ErrAccessDenied
	}

	if member.Role != "admin" {
		return utils.ErrAccessDenied
	}

	// Check if new admin is a member
	newMember, err := s.db.FindOrganizationMember(orgID, newAdminID)
	if err != nil {
		return utils.ErrAccessDenied
	}

	// Update roles
	member.Role = "member"
	newMember.Role = "admin"

	if err := s.db.UpdateOrganizationMember(member); err != nil {
		return err
	}

	return s.db.UpdateOrganizationMember(newMember)
} 