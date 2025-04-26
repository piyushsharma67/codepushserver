package models

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	PublicToken  string    `json:"public_token"`
	PrivateToken string    `json:"private_token,omitempty"`
	CreatedBy    uint      `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type OrganizationInvitation struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	OrganizationID uuid.UUID `json:"organization_id" gorm:"type:uuid;not null"`
	Email          string    `json:"email" gorm:"not null"`
	Role           string    `json:"role" gorm:"not null"`
	Status         string    `json:"status" gorm:"not null;default:'pending'"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type OrganizationRepository interface {
	Create(org *Organization) error
	FindByID(id uuid.UUID) (*Organization, error)
	FindByUserID(userID uint) ([]*Organization, error)
	Update(org *Organization) error
	Delete(id uuid.UUID) error
	AddMember(orgID uuid.UUID, userID uint) error
	RemoveMember(orgID uuid.UUID, userID uint) error
} 