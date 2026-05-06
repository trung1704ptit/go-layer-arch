package repositories

import (
	"github.com/go-layer-arch/internal/models"
	"gorm.io/gorm"
)

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (pr *postRepository) Create(post *models.Post) error {
	result := pr.db.Create(post)
	return result.Error
}

func (pr *postRepository) FindByID(postID string) (*models.Post, error) {
	var post models.Post
	result := pr.db.First(&post, "id = ?", postID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &post, nil
}

func (pr *postRepository) Update(existingPost *models.Post, updates models.Post) error {
	result := pr.db.Model(existingPost).Updates(updates)
	return result.Error
}

func (pr *postRepository) FindAll(limit int, offset int) ([]models.Post, error) {
	var posts []models.Post
	result := pr.db.Limit(limit).Offset(offset).Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}

func (pr *postRepository) DeleteByID(postID string) error {
	result := pr.db.Delete(&models.Post{}, "id = ?", postID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
