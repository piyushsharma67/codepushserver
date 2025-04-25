package database

import (
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
	FindUserByAppID(appID string) (*models.User, error)
	FindUserByToken(token string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error

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