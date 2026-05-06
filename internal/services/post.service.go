package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-layer-arch/internal/models"
	"github.com/go-layer-arch/internal/repositories"
	"gorm.io/gorm"
)

type postService struct {
	postRepository repositories.PostRepository
}

func NewPostService(postRepository repositories.PostRepository) PostService {
	return &postService{postRepository: postRepository}
}

func (ps *postService) Create(payload *models.CreatePostRequest, currentUser models.User) (*models.Post, error) {
	now := time.Now()
	newPost := models.Post{
		Title:     payload.Title,
		Content:   payload.Content,
		Image:     payload.Image,
		User:      currentUser.ID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := ps.postRepository.Create(&newPost)
	if err != nil && strings.Contains(err.Error(), "duplicate key") {
		return nil, ErrPostAlreadyExists
	}
	if err != nil {
		return nil, fmt.Errorf("create post: %w", err)
	}

	return &newPost, nil
}

func (ps *postService) Update(postID string, payload *models.UpdatePost, currentUser models.User) (*models.Post, error) {
	updatedPost, err := ps.postRepository.FindByID(postID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find post by id: %w", err)
		}
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

	if err := ps.postRepository.Update(updatedPost, postToUpdate); err != nil {
		return nil, fmt.Errorf("update post: %w", err)
	}
	return updatedPost, nil
}

func (ps *postService) FindByID(postID string) (*models.Post, error) {
	post, err := ps.postRepository.FindByID(postID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find post by id: %w", err)
		}
		return nil, ErrPostNotFound
	}

	return post, nil
}

func (ps *postService) FindAll(page string, limit string) ([]models.Post, error) {
	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	posts, err := ps.postRepository.FindAll(intLimit, offset)
	if err != nil {
		return nil, fmt.Errorf("find posts: %w", err)
	}

	return posts, nil
}

func (ps *postService) Delete(postID string) error {
	if err := ps.postRepository.DeleteByID(postID); err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("delete post: %w", err)
		}
		return ErrPostNotFound
	}
	return nil
}
