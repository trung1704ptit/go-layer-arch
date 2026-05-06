package repositories

import "github.com/go-layer-arch/internal/models"

type PostRepository interface {
	Create(post *models.Post) error
	FindByID(postID string) (*models.Post, error)
	Update(existingPost *models.Post, updates models.Post) error
	FindAll(limit int, offset int) ([]models.Post, error)
	DeleteByID(postID string) error
}
