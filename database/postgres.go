package database

import (
	"errors"
	"fmt"

	"github.com/piyushsharma67/codepushserver/config"
	"github.com/piyushsharma67/codepushserver/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
	db *gorm.DB
}

func NewPostgresDB(config *config.Config) (Database, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.DBHost,
		config.DBUser,
		config.DBPassword,
		config.DBName,
		config.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}

	return &PostgresDB{db: db}, nil
}

func (p *PostgresDB) Connect() error {
	// Since we're using GORM, the connection is already established when the db is passed in
	return nil
}

func (p *PostgresDB) Close() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (p *PostgresDB) Migrate() error {
	return p.db.AutoMigrate(&models.User{}, &models.App{})
}

func (p *PostgresDB) CreateUser(user *models.User) error {
	return p.db.Create(user).Error
}

func (p *PostgresDB) FindUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := p.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (p *PostgresDB) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := p.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (p *PostgresDB) FindUserByAppID(appID string) (*models.User, error) {
	var user models.User
	if err := p.db.Joins("JOIN apps ON apps.user_id = users.id").
		Where("apps.id = ?", appID).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (p *PostgresDB) FindUserByToken(token string) (*models.User, error) {
	var user models.User
	if err := p.db.Where("token = ?", token).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (p *PostgresDB) UpdateUser(user *models.User) error {
	return p.db.Save(user).Error
}

func (p *PostgresDB) DeleteUser(id uint) error {
	return p.db.Delete(&models.User{}, id).Error
}

func (p *PostgresDB) CreateApp(app *models.App) error {
	return p.db.Create(app).Error
}

func (p *PostgresDB) FindAppByID(id string) (*models.App, error) {
	var app models.App
	if err := p.db.First(&app, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &app, nil
}

func (p *PostgresDB) FindAppsByUserID(userID uint) ([]*models.App, error) {
	var apps []*models.App
	if err := p.db.Where("user_id = ?", userID).Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}

func (p *PostgresDB) UpdateApp(app *models.App) error {
	return p.db.Save(app).Error
}

func (p *PostgresDB) DeleteApp(id string) error {
	return p.db.Delete(&models.App{}, "id = ?", id).Error
} 