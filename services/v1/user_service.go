package v1

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/google/uuid"
	"github.com/piyushsharma67/codepushserver/database"
	"github.com/piyushsharma67/codepushserver/models"
	"github.com/piyushsharma67/codepushserver/utils"
)

type UserService struct {
	db database.Database
}

func NewUserService(db database.Database) *UserService {
	return &UserService{db: db}
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func (s *UserService) GetUserProfile(userID uint) (*models.User, error) {
	return s.db.FindUserByID(userID)
}

func (s *UserService) UpdateUserProfile(userID uint, username, companyName, phoneNumber string) (*models.User, error) {
	user, err := s.db.FindUserByID(userID)
	if err != nil {
		return nil, err
	}

	if username != "" {
		user.Username = username
	}
	if companyName != "" {
		user.CompanyName = companyName
	}
	if phoneNumber != "" {
		user.PhoneNumber = phoneNumber
	}

	if err := s.db.UpdateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) CreateApp(userID uint, name, description string) (*models.App, error) {
	user, err := s.db.FindUserByID(userID)
	if err != nil {
		return nil, err
	}

	app := &models.App{
		ID:          uuid.New().String(),
		UserID:      user.ID,
		Name:        name,
		Description: description,
		Token:       generateRandomString(64),
	}

	if err := s.db.CreateApp(app); err != nil {
		return nil, err
	}

	return app, nil
}

func (s *UserService) GetApp(userID uint, appID string) (*models.App, error) {
	app, err := s.db.FindAppByID(appID)
	if err != nil {
		return nil, err
	}

	if app.UserID != userID {
		return nil, utils.ErrAccessDenied
	}

	return app, nil
}

func (s *UserService) UpdateApp(userID uint, appID string) (*models.App, error) {
	app, err := s.db.FindAppByID(appID)
	if err != nil {
		return nil, err
	}

	if app.UserID != userID {
		return nil, utils.ErrAccessDenied
	}

	app.Token = generateRandomString(64)
	if err := s.db.UpdateApp(app); err != nil {
		return nil, err
	}

	return app, nil
}

func (s *UserService) DeleteApp(userID uint, appID string) error {
	app, err := s.db.FindAppByID(appID)
	if err != nil {
		return err
	}

	if app.UserID != userID {
		return utils.ErrAccessDenied
	}

	return s.db.DeleteApp(app.ID)
}

func (s *UserService) GetAllApps(userID uint) ([]*models.App, error) {
	return s.db.FindAppsByUserID(userID)
} 