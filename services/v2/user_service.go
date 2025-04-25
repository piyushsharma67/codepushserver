package v2

import (
	"github.com/piyushsharma67/codepushserver/errors"
	"github.com/piyushsharma67/codepushserver/models"
	store "github.com/piyushsharma67/codepushserver/store/v2"
)

type UserService struct {
	store store.Store
}

func NewUserService(store store.Store) *UserService {
	return &UserService{store: store}
}

type UpdateProfileRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	CompanyName string `json:"company_name"`
	PhoneNumber string `json:"phone_number"`
}

type CreateAppRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Platform    string `json:"platform" binding:"required,oneof=ios android"`
}

type UpdateAppRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Platform    string `json:"platform" binding:"omitempty,oneof=ios android"`
}

func (s *UserService) GetProfile(userID uint) (*models.User, error) {
	user, err := s.store.FindUserByID(userID)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return user, nil
}

func (s *UserService) UpdateProfile(userID uint, req *UpdateProfileRequest) (*models.User, error) {
	user, err := s.store.FindUserByID(userID)
	if err != nil {
		return nil, errors.ErrNotFound
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.CompanyName != "" {
		user.CompanyName = req.CompanyName
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = req.PhoneNumber
	}

	if err := s.store.UpdateUser(user); err != nil {
		return nil, errors.ErrInternal
	}

	return user, nil
}

func (s *UserService) CreateApp(userID uint, name, description, platform string) (*models.App, error) {
	app := &models.App{
		UserID:      userID,
		Name:        name,
		Description: description,
		Platform:    platform,
	}

	if err := s.store.CreateApp(app); err != nil {
		return nil, errors.ErrInternal
	}

	return app, nil
}

func (s *UserService) GetApp(userID uint, appID string) (*models.App, error) {
	app, err := s.store.FindAppByID(appID)
	if err != nil {
		return nil, errors.ErrNotFound
	}

	if app.UserID != userID {
		return nil, errors.ErrUnauthorized
	}

	return app, nil
}

func (s *UserService) UpdateApp(userID uint, appID string, name, description, platform string) (*models.App, error) {
	app, err := s.store.FindAppByID(appID)
	if err != nil {
		return nil, errors.ErrNotFound
	}

	if app.UserID != userID {
		return nil, errors.ErrUnauthorized
	}

	if name != "" {
		app.Name = name
	}
	if description != "" {
		app.Description = description
	}
	if platform != "" {
		app.Platform = platform
	}

	if err := s.store.UpdateApp(app); err != nil {
		return nil, errors.ErrInternal
	}

	return app, nil
}

func (s *UserService) DeleteApp(userID uint, appID string) error {
	app, err := s.store.FindAppByID(appID)
	if err != nil {
		return errors.ErrNotFound
	}

	if app.UserID != userID {
		return errors.ErrUnauthorized
	}

	if err := s.store.DeleteApp(appID); err != nil {
		return errors.ErrInternal
	}

	return nil
}

func (s *UserService) ListApps(userID uint) ([]*models.App, error) {
	apps, err := s.store.FindAppsByUserID(userID)
	if err != nil {
		return nil, errors.ErrInternal
	}
	return apps, nil
} 