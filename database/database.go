package database

import (
	"github.com/google/uuid"
	"github.com/piyushsharma67/codepushserver/config"
	"github.com/piyushsharma67/codepushserver/models"
)

type Database interface {
	Connect() error
	Close() error
	Migrate() error

	// User methods
	CreateUser(user *models.User) error
	FindUserByID(id uint) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error

	// Organization methods
	CreateOrganization(org *models.Organization) error
	FindOrganizationByID(id uuid.UUID) (*models.Organization, error)
	FindOrganizationsByUserID(userID uint) ([]*models.Organization, error)
	DeleteOrganization(id uuid.UUID) error

	// Organization member methods
	CreateOrganizationMember(member *models.OrganizationMember) error
	FindOrganizationMember(orgID uuid.UUID, userID uint) (*models.OrganizationMember, error)
	UpdateOrganizationMember(member *models.OrganizationMember) error

	// Organization invitation methods
	CreateOrganizationInvitation(invitation *models.OrganizationInvitation) error
	FindOrganizationInvitationByID(id uuid.UUID) (*models.OrganizationInvitation, error)
	FindPendingInvitationsByEmail(email string) ([]*models.OrganizationInvitation, error)
	UpdateOrganizationInvitation(invitation *models.OrganizationInvitation) error

	// App methods
	CreateApp(app *models.App) error
	FindAppByID(id string) (*models.App, error)
	FindAppsByUserID(userID uint) ([]*models.App, error)
	UpdateApp(app *models.App) error
	DeleteApp(id string) error
}

// NewDatabase creates a new database instance based on the configuration
func NewDatabase(config *config.Config) (Database, error) {
	switch config.DBType {
	case "postgres":
		return NewPostgresDB(config)
	case "mysql":
		return NewMySQLDB(config)
	default:
		return NewPostgresDB(config)
	}
} 