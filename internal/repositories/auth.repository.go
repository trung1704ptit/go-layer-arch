package repositories

import (
	"strings"
	"time"

	"github.com/go-layer-arch/internal/models"
	"gorm.io/gorm"
)

type AuthRepository interface {
	CreateUser(payload *models.SignUpInput, hashedPassword string) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)
	FindUserByID(userID string) (*models.User, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (ar *authRepository) CreateUser(payload *models.SignUpInput, hashedPassword string) (*models.User, error) {
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

	result := ar.db.Create(&newUser)
	if result.Error != nil {
		return nil, result.Error
	}

	return &newUser, nil
}

func (ar *authRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := ar.db.First(&user, "email = ?", strings.ToLower(email))
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (ar *authRepository) FindUserByID(userID string) (*models.User, error) {
	var user models.User
	result := ar.db.First(&user, "id = ?", userID)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
