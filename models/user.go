package models

import (
	"time"
)

type User struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Username    string    `json:"username" gorm:"unique;not null"`
	Email       string    `json:"email" gorm:"unique;not null"`
	Password    string    `json:"-" gorm:"not null"`
	CompanyName string    `json:"company_name"`
	PhoneNumber string    `json:"phone_number"`
	AppID       string    `json:"app_id" gorm:"unique;not null"`
	Token       string    `json:"token" gorm:"unique;not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByAppID(appID string) (*User, error)
	FindByToken(token string) (*User, error)
	Update(user *User) error
	Delete(id uint) error
} 