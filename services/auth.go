package services

import (
	"taskive/config"
	"taskive/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

type RegisterInput struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (s *AuthService) Register(input RegisterInput) (*models.User, error) {
	// Check if email exists using a custom query to avoid GORM logging
	var count int64
	if err := s.db.Model(&models.User{}).Where("email = ?", input.Email).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, models.ErrEmailExists
	}

	user := &models.User{
		Name:  input.Name,
		Email: input.Email,
	}

	if err := user.SetPassword(input.Password); err != nil {
		return nil, err
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(input LoginInput) (string, error) {
	var user models.User
	if err := s.db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", models.ErrUserNotFound
		}
		return "", err
	}

	if !user.CheckPassword(input.Password) {
		return "", models.ErrInvalidPassword
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
} 