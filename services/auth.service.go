package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-layer-arch/initializers"
	"github.com/go-layer-arch/models"
	"github.com/go-layer-arch/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

func NewAuthService(db *gorm.DB) AuthService {
	return AuthService{DB: db}
}

func (as *AuthService) SignUp(payload *models.SignUpInput) (*models.UserResponse, error) {
	if payload.Password != payload.PasswordConfirm {
		return nil, ErrPasswordsDoNotMatch
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	newUser := models.User{
		Name:      payload.Name,
		Email:     strings.ToLower(payload.Email),
		Password:  hashedPassword,
		Role:      "user",
		Verified:  true,
		Photo:     payload.Photo,
		Provider:  "local",
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := as.DB.Create(&newUser)
	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		return nil, ErrUserAlreadyExists
	}
	if result.Error != nil {
		return nil, fmt.Errorf("create user: %w", result.Error)
	}

	return &models.UserResponse{
		ID:        newUser.ID,
		Name:      newUser.Name,
		Email:     newUser.Email,
		Photo:     newUser.Photo,
		Role:      newUser.Role,
		Provider:  newUser.Provider,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}, nil
}

func (as *AuthService) SignIn(payload *models.SignInInput) (string, string, error) {
	var user models.User
	result := as.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		return "", "", ErrInvalidCredentials
	}

	if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
		return "", "", ErrInvalidCredentials
	}

	config, _ := initializers.LoadConfig(".")
	accessToken, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := utils.CreateToken(config.RefreshTokenExpiresIn, user.ID, config.RefreshTokenPrivateKey)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (as *AuthService) RefreshAccessToken(refreshToken string) (string, error) {
	config, _ := initializers.LoadConfig(".")

	sub, err := utils.ValidateToken(refreshToken, config.RefreshTokenPublicKey)
	if err != nil {
		return "", err
	}

	var user models.User
	result := as.DB.First(&user, "id = ?", fmt.Sprint(sub))
	if result.Error != nil {
		return "", ErrRefreshUserNotFound
	}

	accessToken, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
