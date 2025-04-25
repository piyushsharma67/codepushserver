package v2

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/piyushsharma67/codepushserver/config"
	"github.com/piyushsharma67/codepushserver/errors"
	"github.com/piyushsharma67/codepushserver/models"
	store "github.com/piyushsharma67/codepushserver/store/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	store store.Store
}

func NewAuthService(store store.Store) *AuthService {
	return &AuthService{store: store}
}

type RegisterRequest struct {
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	CompanyName string `json:"company_name"`
	PhoneNumber string `json:"phone_number"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      struct {
		ID          uint      `json:"id"`
		Email       string    `json:"email"`
		Username    string    `json:"username"`
		CompanyName string    `json:"company_name"`
		AppID       string    `json:"app_id"`
		Token       string    `json:"token"`
		CreatedAt   time.Time `json:"created_at"`
	} `json:"user"`
}

func (s *AuthService) Register(req *RegisterRequest) (*models.User, error) {
	// Check if user already exists
	existingUser, _ := s.store.FindUserByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrInternal
	}

	// Generate app ID and token
	appID := generateRandomString(32)
	token := generateRandomString(64)

	// Create user
	user := &models.User{
		Username:    req.Username,
		Email:       req.Email,
		Password:    string(hashedPassword),
		AppID:       appID,
		Token:       token,
		CompanyName: req.CompanyName,
		PhoneNumber: req.PhoneNumber,
	}

	if err := s.store.CreateUser(user); err != nil {
		return nil, errors.ErrInternal
	}

	return user, nil
}

func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	// Find user
	user, err := s.store.FindUserByEmail(req.Email)
	if err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	// Generate JWT token
	expiresAt := time.Now().Add(time.Hour * 24)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     expiresAt.Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.LoadConfig().JWTSecret))
	if err != nil {
		return nil, errors.ErrInternal
	}

	response := &LoginResponse{
		Token:     tokenString,
		ExpiresAt: expiresAt,
	}
	response.User.ID = user.ID
	response.User.Email = user.Email
	response.User.Username = user.Username
	response.User.CompanyName = user.CompanyName
	response.User.AppID = user.AppID
	response.User.Token = user.Token
	response.User.CreatedAt = user.CreatedAt

	return response, nil
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(result)
} 