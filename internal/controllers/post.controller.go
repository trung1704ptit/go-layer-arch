package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-layer-arch/internal/models"
	"github.com/go-layer-arch/internal/services"
	"github.com/go-layer-arch/pkg/shared"
)

type PostController struct {
	postService services.PostService
}

func NewPostController(postService services.PostService) PostController {
	return PostController{postService: postService}
}

func (pc *PostController) CreatePost(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload *models.CreatePostRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		shared.WriteError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	newPost, err := pc.postService.Create(payload, currentUser)
	if err != nil {
		if errors.Is(err, services.ErrPostAlreadyExists) {
			shared.WriteError(ctx, http.StatusConflict, "Post with that title already exists")
		} else {
			shared.WriteError(ctx, http.StatusBadGateway, "Something bad happened")
		}
		return
	}

	shared.WriteSuccess(ctx, http.StatusCreated, newPost)
}

func (pc *PostController) UpdatePost(ctx *gin.Context) {
	postId := ctx.Param("postId")
	currentUser := ctx.MustGet("currentUser").(models.User)

	var payload *models.UpdatePost
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		shared.WriteError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	updatedPost, err := pc.postService.Update(postId, payload, currentUser)
	if err != nil {
		if errors.Is(err, services.ErrPostNotFound) {
			shared.WriteError(ctx, http.StatusNotFound, "No post with that title exists")
		} else {
			shared.WriteError(ctx, http.StatusBadGateway, "Something bad happened")
		}
		return
	}

	shared.WriteSuccess(ctx, http.StatusOK, updatedPost)
}

func (pc *PostController) FindPostById(ctx *gin.Context) {
	postId := ctx.Param("postId")

	post, err := pc.postService.FindByID(postId)
	if err != nil {
		if errors.Is(err, services.ErrPostNotFound) {
			shared.WriteError(ctx, http.StatusNotFound, "No post with that title exists")
		} else {
			shared.WriteError(ctx, http.StatusBadGateway, "Something bad happened")
		}
		return
	}

	shared.WriteSuccess(ctx, http.StatusOK, post)
}

func (pc *PostController) FindPosts(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	posts, err := pc.postService.FindAll(page, limit)
	if err != nil {
		shared.WriteError(ctx, http.StatusBadGateway, "Something bad happened")
		return
	}

	shared.WriteSuccess(ctx, http.StatusOK, gin.H{"results": len(posts), "posts": posts})
}

func (pc *PostController) DeletePost(ctx *gin.Context) {
	postId := ctx.Param("postId")

	if err := pc.postService.Delete(postId); err != nil {
		if errors.Is(err, services.ErrPostNotFound) {
			shared.WriteError(ctx, http.StatusNotFound, "No post with that title exists")
		} else {
			shared.WriteError(ctx, http.StatusBadGateway, "Something bad happened")
		}
		return
	}

	shared.WriteSuccess(ctx, http.StatusOK, gin.H{"deleted": true})
}
