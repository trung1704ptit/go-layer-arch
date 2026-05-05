package services

import "github.com/go-layer-arch/models"

type UserService struct{}

func NewUserService() UserService {
	return UserService{}
}

func (us *UserService) BuildUserResponse(user models.User) *models.UserResponse {
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
