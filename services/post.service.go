package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-layer-arch/models"
	"github.com/go-layer-arch/repositories"
)

type PostService struct {
	postRepository repositories.PostRepository
}

func NewPostService(postRepository repositories.PostRepository) PostService {
	return PostService{postRepository: postRepository}
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

	err := ps.postRepository.Create(&newPost)
	if err != nil && strings.Contains(err.Error(), "duplicate key") {
		return nil, ErrPostAlreadyExists
	}
	if err != nil {
		return nil, fmt.Errorf("create post: %w", err)
	}

	return &newPost, nil
}

func (ps *PostService) Update(postID string, payload *models.UpdatePost, currentUser models.User) (*models.Post, error) {
	updatedPost, err := ps.postRepository.FindByID(postID)
	if err != nil {
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

func (ps *PostService) FindByID(postID string) (*models.Post, error) {
	post, err := ps.postRepository.FindByID(postID)
	if err != nil {
		return nil, ErrPostNotFound
	}

	return post, nil
}

func (ps *PostService) FindAll(page string, limit string) ([]models.Post, error) {
	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	posts, err := ps.postRepository.FindAll(intLimit, offset)
	if err != nil {
		return nil, fmt.Errorf("find posts: %w", err)
	}

	return posts, nil
}

func (ps *PostService) Delete(postID string) error {
	if err := ps.postRepository.DeleteByID(postID); err != nil {
		return ErrPostNotFound
	}
	return nil
}
