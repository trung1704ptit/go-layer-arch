package services

import "github.com/go-layer-arch/internal/models"

type PostService interface {
	Create(payload *models.CreatePostRequest, currentUser models.User) (*models.Post, error)
	Update(postID string, payload *models.UpdatePost, currentUser models.User) (*models.Post, error)
	FindByID(postID string) (*models.Post, error)
	FindAll(page string, limit string) ([]models.Post, error)
	Delete(postID string) error
}
