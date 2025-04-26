package models

import (
	"time"
)

type User struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Username    string    `json:"username" gorm:"not null"`
	Email       string    `json:"email" gorm:"unique;not null"`
	Password    string    `json:"-" gorm:"not null"`
	CompanyName string    `json:"company_name"`
	PhoneNumber string    `json:"phone_number"`
	Apps        []App     `json:"apps" gorm:"foreignKey:UserID"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uint) error
}

type AppRepository interface {
	Create(app *App) error
	FindByAppID(appID string) (*App, error)
	FindByUserID(userID uint) ([]App, error)
	Update(app *App) error
	Delete(id uint) error
} 