package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-layer-arch/internal/models"
	"github.com/go-layer-arch/internal/services"
	"github.com/go-layer-arch/pkg/shared"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return UserController{userService: userService}
}

func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	userResponse := uc.userService.BuildUserResponse(currentUser)

	shared.WriteSuccess(ctx, http.StatusOK, gin.H{"user": userResponse})
}
