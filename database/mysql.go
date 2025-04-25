package database

import (
	"errors"
	"fmt"

	"github.com/piyushsharma67/codepushserver/config"
	"github.com/piyushsharma67/codepushserver/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLDB struct {
	db *gorm.DB
}

func (m *MySQLDB) Connect() error {
	// Since we're using GORM, the connection is already established when the db is passed in
	return nil
}

func (m *MySQLDB) Close() error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (m *MySQLDB) Migrate() error {
	return m.db.AutoMigrate(&models.User{}, &models.App{})
}

func (m *MySQLDB) CreateUser(user *models.User) error {
	return m.db.Create(user).Error
}

func (m *MySQLDB) FindUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := m.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (m *MySQLDB) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := m.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (m *MySQLDB) FindUserByAppID(appID string) (*models.User, error) {
	var user models.User
	if err := m.db.Joins("JOIN apps ON apps.user_id = users.id").
		Where("apps.id = ?", appID).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (m *MySQLDB) FindUserByToken(token string) (*models.User, error) {
	var user models.User
	if err := m.db.Where("token = ?", token).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (m *MySQLDB) UpdateUser(user *models.User) error {
	return m.db.Save(user).Error
}

func (m *MySQLDB) DeleteUser(id uint) error {
	return m.db.Delete(&models.User{}, id).Error
}

func (m *MySQLDB) CreateApp(app *models.App) error {
	return m.db.Create(app).Error
}

func (m *MySQLDB) FindAppByID(id string) (*models.App, error) {
	var app models.App
	if err := m.db.First(&app, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &app, nil
}

func (m *MySQLDB) FindAppsByUserID(userID uint) ([]*models.App, error) {
	var apps []*models.App
	if err := m.db.Where("user_id = ?", userID).Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}

func (m *MySQLDB) UpdateApp(app *models.App) error {
	return m.db.Save(app).Error
}

func (m *MySQLDB) DeleteApp(id string) error {
	return m.db.Delete(&models.App{}, "id = ?", id).Error
}

func NewMySQLDB(config *config.Config) (Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %v", err)
	}

	return &MySQLDB{db: db}, nil
} 