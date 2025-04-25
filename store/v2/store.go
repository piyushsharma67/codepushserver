package v2

import (
	"github.com/piyushsharma67/codepushserver/models"
)

type Store interface {
	// User operations
	CreateUser(user *models.User) error
	FindUserByID(id uint) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error

	// App operations
	CreateApp(app *models.App) error
	FindAppByID(id string) (*models.App, error)
	FindAppsByUserID(userID uint) ([]*models.App, error)
	UpdateApp(app *models.App) error
	DeleteApp(id string) error
} 