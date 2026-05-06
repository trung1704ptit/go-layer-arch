package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-layer-arch/internal/initializers"
	"github.com/go-layer-arch/internal/models"
	"github.com/go-layer-arch/internal/repositories"
	"github.com/go-layer-arch/pkg/utils"
	"gorm.io/gorm"
)

type authService struct {
	config         initializers.Config
	authRepository repositories.AuthRepository
}

func NewAuthService(config initializers.Config, authRepository repositories.AuthRepository) AuthService {
	return &authService{
		config:         config,
		authRepository: authRepository,
	}
}

func (as *authService) SignUp(payload *models.SignUpInput) (*models.UserResponse, error) {
	if payload.Password != payload.PasswordConfirm {
		return nil, ErrPasswordsDoNotMatch
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	newUser, err := as.authRepository.CreateUser(payload, hashedPassword)
	if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique") {
		return nil, ErrUserAlreadyExists
	}
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
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

func (as *authService) SignIn(payload *models.SignInInput) (string, string, error) {
	user, err := as.authRepository.FindUserByEmail(payload.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", ErrInvalidCredentials
		}
		return "", "", fmt.Errorf("find user by email: %w", err)
	}

	if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
		return "", "", ErrInvalidCredentials
	}

	accessToken, err := utils.CreateToken(as.config.AccessTokenExpiresIn, user.ID, as.config.AccessTokenPrivateKey)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := utils.CreateToken(as.config.RefreshTokenExpiresIn, user.ID, as.config.RefreshTokenPrivateKey)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (as *authService) RefreshAccessToken(refreshToken string) (string, error) {
	sub, err := utils.ValidateToken(refreshToken, as.config.RefreshTokenPublicKey)
	if err != nil {
		return "", ErrRefreshTokenInvalid
	}

	user, err := as.authRepository.FindUserByID(fmt.Sprint(sub))
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("find refresh user: %w", err)
		}
		return "", ErrRefreshUserNotFound
	}

	accessToken, err := utils.CreateToken(as.config.AccessTokenExpiresIn, user.ID, as.config.AccessTokenPrivateKey)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
