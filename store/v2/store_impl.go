package v2

import (
	"github.com/piyushsharma67/codepushserver/database"
	"github.com/piyushsharma67/codepushserver/models"
)

type storeImpl struct {
	db database.Database
}

func NewStore(db database.Database) Store {
	return &storeImpl{db: db}
}

func (s *storeImpl) CreateUser(user *models.User) error {
	return s.db.CreateUser(user)
}

func (s *storeImpl) FindUserByID(id uint) (*models.User, error) {
	return s.db.FindUserByID(id)
}

func (s *storeImpl) FindUserByEmail(email string) (*models.User, error) {
	return s.db.FindUserByEmail(email)
}

func (s *storeImpl) UpdateUser(user *models.User) error {
	return s.db.UpdateUser(user)
}

func (s *storeImpl) CreateApp(app *models.App) error {
	return s.db.CreateApp(app)
}

func (s *storeImpl) FindAppByID(id string) (*models.App, error) {
	return s.db.FindAppByID(id)
}

func (s *storeImpl) FindAppsByUserID(userID uint) ([]*models.App, error) {
	return s.db.FindAppsByUserID(userID)
}

func (s *storeImpl) UpdateApp(app *models.App) error {
	return s.db.UpdateApp(app)
}

func (s *storeImpl) DeleteApp(id string) error {
	return s.db.DeleteApp(id)
} 