package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-layer-arch/models"
	"gorm.io/gorm"
)

type PostService struct {
	DB *gorm.DB
}

func NewPostService(db *gorm.DB) PostService {
	return PostService{DB: db}
}

func (ps *PostService) Create(payload *models.CreatePostRequest, currentUser models.User) (*models.Post, error) {
	now := time.Now()
	newPost := models.Post{
		Title:     payload.Title,
		Content:   payload.Content,
		Image:     payload.Image,
		User:      currentUser.ID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := ps.DB.Create(&newPost)
	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key") {
		return nil, ErrPostAlreadyExists
	}
	if result.Error != nil {
		return nil, fmt.Errorf("create post: %w", result.Error)
	}

	return &newPost, nil
}

func (ps *PostService) Update(postID string, payload *models.UpdatePost, currentUser models.User) (*models.Post, error) {
	var updatedPost models.Post
	result := ps.DB.First(&updatedPost, "id = ?", postID)
	if result.Error != nil {
		return nil, ErrPostNotFound
	}

	postToUpdate := models.Post{
		Title:     payload.Title,
		Content:   payload.Content,
		Image:     payload.Image,
		User:      currentUser.ID,
		CreatedAt: updatedPost.CreatedAt,
		UpdatedAt: time.Now(),
	}

	ps.DB.Model(&updatedPost).Updates(postToUpdate)
	return &updatedPost, nil
}

func (ps *PostService) FindByID(postID string) (*models.Post, error) {
	var post models.Post
	result := ps.DB.First(&post, "id = ?", postID)
	if result.Error != nil {
		return nil, ErrPostNotFound
	}

	return &post, nil
}

func (ps *PostService) FindAll(page string, limit string) ([]models.Post, error) {
	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var posts []models.Post
	result := ps.DB.Limit(intLimit).Offset(offset).Find(&posts)
	if result.Error != nil {
		return nil, fmt.Errorf("find posts: %w", result.Error)
	}

	return posts, nil
}

func (ps *PostService) Delete(postID string) error {
	result := ps.DB.Delete(&models.Post{}, "id = ?", postID)
	if result.Error != nil {
		return ErrPostNotFound
	}
	return nil
}
