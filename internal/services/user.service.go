package services

import "github.com/go-layer-arch/internal/models"

type UserService interface {
	BuildUserResponse(user models.User) *models.UserResponse
}

type userService struct{}

func NewUserService() UserService {
	return &userService{}
}

func (us *userService) BuildUserResponse(user models.User) *models.UserResponse {
	return &models.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Photo:     user.Photo,
		Role:      user.Role,
		Provider:  user.Provider,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
