package repositories

import "github.com/go-layer-arch/internal/models"

type AuthRepository interface {
	CreateUser(payload *models.SignUpInput, hashedPassword string) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)
	FindUserByID(userID string) (*models.User, error)
}
