package database

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/piyushsharma67/codepushserver/config"
	"github.com/piyushsharma67/codepushserver/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
	db *gorm.DB
}

func NewPostgresDB(config *config.Config) (*PostgresDB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &PostgresDB{db: db}, nil
}

func (d *PostgresDB) Connect() error {
	return nil
}

func (d *PostgresDB) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *PostgresDB) Migrate() error {
	// Enable UUID extension
	if err := d.db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return err
	}

	// Drop existing tables to ensure clean migration
	if err := d.db.Migrator().DropTable(
		&models.Organization{},
		&models.OrganizationMember{},
		&models.OrganizationInvitation{},
	); err != nil {
		return err
	}

	// Create tables with new schema
	return d.db.AutoMigrate(
		&models.User{},
		&models.App{},
		&models.Organization{},
		&models.OrganizationMember{},
		&models.OrganizationInvitation{},
	)
}

// User methods
func (d *PostgresDB) CreateUser(user *models.User) error {
	return d.db.Create(user).Error
}

func (d *PostgresDB) FindUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := d.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *PostgresDB) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := d.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *PostgresDB) UpdateUser(user *models.User) error {
	return d.db.Save(user).Error
}

// App methods
func (d *PostgresDB) CreateApp(app *models.App) error {
	return d.db.Create(app).Error
}

func (d *PostgresDB) FindAppByID(id string) (*models.App, error) {
	var app models.App
	if err := d.db.First(&app, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

func (d *PostgresDB) FindAppsByUserID(userID uint) ([]*models.App, error) {
	var apps []*models.App
	if err := d.db.Where("user_id = ?", userID).Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}

func (d *PostgresDB) UpdateApp(app *models.App) error {
	return d.db.Save(app).Error
}

func (d *PostgresDB) DeleteApp(id string) error {
	return d.db.Delete(&models.App{}, "id = ?", id).Error
}

// Organization methods
func (d *PostgresDB) CreateOrganization(org *models.Organization) error {
	return d.db.Create(org).Error
}

func (d *PostgresDB) FindOrganizationByID(id uuid.UUID) (*models.Organization, error) {
	var org models.Organization
	if err := d.db.First(&org, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (d *PostgresDB) FindOrganizationsByUserID(userID uint) ([]*models.Organization, error) {
	var orgs []*models.Organization
	if err := d.db.Joins("JOIN organization_members ON organizations.id::uuid = organization_members.organization_id::uuid").
		Where("organization_members.user_id = ?", userID).
		Find(&orgs).Error; err != nil {
		return nil, err
	}
	return orgs, nil
}

func (d *PostgresDB) DeleteOrganization(id uuid.UUID) error {
	return d.db.Delete(&models.Organization{}, "id = ?", id).Error
}

// Organization member methods
func (d *PostgresDB) CreateOrganizationMember(member *models.OrganizationMember) error {
	return d.db.Create(member).Error
}

func (d *PostgresDB) FindOrganizationMember(orgID uuid.UUID, userID uint) (*models.OrganizationMember, error) {
	var member models.OrganizationMember
	if err := d.db.Where("organization_id = ? AND user_id = ?", orgID, userID).First(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

func (d *PostgresDB) UpdateOrganizationMember(member *models.OrganizationMember) error {
	return d.db.Save(member).Error
}

// Organization invitation methods
func (d *PostgresDB) CreateOrganizationInvitation(invitation *models.OrganizationInvitation) error {
	return d.db.Create(invitation).Error
}

func (d *PostgresDB) FindOrganizationInvitationByID(id uuid.UUID) (*models.OrganizationInvitation, error) {
	var invitation models.OrganizationInvitation
	if err := d.db.First(&invitation, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (d *PostgresDB) FindPendingInvitationsByEmail(email string) ([]*models.OrganizationInvitation, error) {
	var invitations []*models.OrganizationInvitation
	if err := d.db.Where("email = ? AND status = ?", email, "pending").Find(&invitations).Error; err != nil {
		return nil, err
	}
	return invitations, nil
}

func (d *PostgresDB) UpdateOrganizationInvitation(invitation *models.OrganizationInvitation) error {
	return d.db.Save(invitation).Error
} 