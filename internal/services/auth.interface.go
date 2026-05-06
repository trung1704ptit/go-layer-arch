package services

import "github.com/go-layer-arch/internal/models"

type AuthService interface {
	SignUp(payload *models.SignUpInput) (*models.UserResponse, error)
	SignIn(payload *models.SignInInput) (string, string, error)
	RefreshAccessToken(refreshToken string) (string, error)
}
