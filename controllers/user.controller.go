package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-layer-arch/models"
	"github.com/go-layer-arch/services"
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

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": userResponse}})
}
