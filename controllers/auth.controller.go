package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-layer-arch/initializers"
	"github.com/go-layer-arch/models"
	"github.com/go-layer-arch/services"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) AuthController {
	return AuthController{authService: authService}
}

// SignUp User
func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var payload *models.SignUpInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	userResponse, err := ac.authService.SignUp(payload)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrPasswordsDoNotMatch):
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		case errors.Is(err, services.ErrUserAlreadyExists):
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that email already exists"})
		default:
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Something bad happened"})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{"user": userResponse}})
}

func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var payload *models.SignInInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	accessToken, refreshToken, err := ac.authService.SignIn(payload)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	config, _ := initializers.LoadConfig(".")

	ctx.SetCookie("access_token", accessToken, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

// Refresh Access Token
func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {
	cookie, err := ctx.Cookie("refresh_token")

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "could not refresh access token"})
		return
	}

	accessToken, err := ac.authService.RefreshAccessToken(cookie)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	config, _ := initializers.LoadConfig(".")
	ctx.SetCookie("access_token", accessToken, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
