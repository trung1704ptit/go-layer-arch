package services

import (
	"fmt"
	"strings"

	"github.com/go-layer-arch/initializers"
	"github.com/go-layer-arch/models"
	"github.com/go-layer-arch/repositories"
	"github.com/go-layer-arch/utils"
)

type AuthService struct {
	authRepository repositories.AuthRepository
}

func NewAuthService(authRepository repositories.AuthRepository) AuthService {
	return AuthService{authRepository: authRepository}
}

func (as *AuthService) SignUp(payload *models.SignUpInput) (*models.UserResponse, error) {
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

func (as *AuthService) SignIn(payload *models.SignInInput) (string, string, error) {
	user, err := as.authRepository.FindUserByEmail(payload.Email)
	if err != nil {
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

	user, err := as.authRepository.FindUserByID(fmt.Sprint(sub))
	if err != nil {
		return "", ErrRefreshUserNotFound
	}

	accessToken, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
