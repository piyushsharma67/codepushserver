package database

import (
	"errors"

	"github.com/piyushsharma67/codepushserver/models"
	"gorm.io/gorm"
)

type databaseImpl struct {
	db *gorm.DB
}

func (d *databaseImpl) Connect() error {
	// Since we're using GORM, the connection is already established when the db is passed in
	return nil
}

func (d *databaseImpl) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *databaseImpl) Migrate() error {
	return d.db.AutoMigrate(&models.User{}, &models.App{})
}

func (d *databaseImpl) CreateUser(user *models.User) error {
	return d.db.Create(user).Error
}

func (d *databaseImpl) FindUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := d.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (d *databaseImpl) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := d.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (d *databaseImpl) FindUserByAppID(appID string) (*models.User, error) {
	var user models.User
	if err := d.db.Joins("JOIN apps ON apps.user_id = users.id").
		Where("apps.id = ?", appID).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (d *databaseImpl) FindUserByToken(token string) (*models.User, error) {
	var user models.User
	if err := d.db.Where("token = ?", token).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (d *databaseImpl) UpdateUser(user *models.User) error {
	return d.db.Save(user).Error
}

func (d *databaseImpl) DeleteUser(id uint) error {
	return d.db.Delete(&models.User{}, id).Error
}

func (d *databaseImpl) CreateApp(app *models.App) error {
	return d.db.Create(app).Error
}

func (d *databaseImpl) FindAppByID(id string) (*models.App, error) {
	var app models.App
	if err := d.db.First(&app, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &app, nil
}

func (d *databaseImpl) FindAppsByUserID(userID uint) ([]*models.App, error) {
	var apps []*models.App
	if err := d.db.Where("user_id = ?", userID).Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}

func (d *databaseImpl) UpdateApp(app *models.App) error {
	return d.db.Save(app).Error
}

func (d *databaseImpl) DeleteApp(id string) error {
	return d.db.Delete(&models.App{}, "id = ?", id).Error
} 