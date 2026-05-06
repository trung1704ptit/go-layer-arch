package repositories

import (
	"github.com/go-layer-arch/models"
	"gorm.io/gorm"
)

type PostRepository struct {
	DB *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return PostRepository{DB: db}
}

func (pr *PostRepository) Create(post *models.Post) error {
	result := pr.DB.Create(post)
	return result.Error
}

func (pr *PostRepository) FindByID(postID string) (*models.Post, error) {
	var post models.Post
	result := pr.DB.First(&post, "id = ?", postID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &post, nil
}

func (pr *PostRepository) Update(existingPost *models.Post, updates models.Post) error {
	result := pr.DB.Model(existingPost).Updates(updates)
	return result.Error
}

func (pr *PostRepository) FindAll(limit int, offset int) ([]models.Post, error) {
	var posts []models.Post
	result := pr.DB.Limit(limit).Offset(offset).Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}

func (pr *PostRepository) DeleteByID(postID string) error {
	result := pr.DB.Delete(&models.Post{}, "id = ?", postID)
	return result.Error
}
