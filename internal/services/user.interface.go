package services

import "github.com/go-layer-arch/internal/models"

type UserService interface {
	BuildUserResponse(user models.User) *models.UserResponse
}
