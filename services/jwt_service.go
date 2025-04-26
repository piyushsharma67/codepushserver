package services

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/piyushsharma67/codepushserver/config"
)

type JWTService struct {
	secretKey []byte
}

func NewJWTService() *JWTService {
	cfg := config.NewConfig()
	return &JWTService{
		secretKey: []byte(cfg.JWTKey),
	}
}

func (s *JWTService) GenerateToken(userID uint, email string) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Hour * 24) // 24 hours

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func (s *JWTService) ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, jwt.ErrSignatureInvalid
} 