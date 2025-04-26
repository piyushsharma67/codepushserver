package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	RoleAdmin     = "admin"
	RoleDeveloper = "developer"
)

var (
	ErrInvalidRole = errors.New("invalid role")
	ErrNotAdmin    = errors.New("user is not an admin")
)

type OrganizationMember struct {
	OrganizationID uuid.UUID      `json:"organization_id" gorm:"primaryKey"`
	UserID         uint      `json:"user_id" gorm:"primaryKey"`
	Role           string    `json:"role" gorm:"not null;default:'developer'"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (m *OrganizationMember) IsAdmin() bool {
	return m.Role == RoleAdmin
}

func (m *OrganizationMember) IsDeveloper() bool {
	return m.Role == RoleDeveloper
}

func (m *OrganizationMember) ValidateRole() error {
	switch m.Role {
	case RoleAdmin, RoleDeveloper:
		return nil
	default:
		return ErrInvalidRole
	}
} 