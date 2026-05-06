package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-layer-arch/internal/models"
	"github.com/go-layer-arch/internal/services"
	"github.com/go-layer-arch/pkg/shared"
)

type AuthController struct {
	authService        services.AuthService
	accessTokenMaxAge  int
	refreshTokenMaxAge int
}

func NewAuthController(authService services.AuthService, accessTokenMaxAge int, refreshTokenMaxAge int) AuthController {
	return AuthController{
		authService:        authService,
		accessTokenMaxAge:  accessTokenMaxAge,
		refreshTokenMaxAge: refreshTokenMaxAge,
	}
}

// SignUp User
func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var payload *models.SignUpInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		shared.WriteError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userResponse, err := ac.authService.SignUp(payload)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrPasswordsDoNotMatch):
			shared.WriteError(ctx, http.StatusBadRequest, "Passwords do not match")
		case errors.Is(err, services.ErrUserAlreadyExists):
			shared.WriteError(ctx, http.StatusConflict, "User with that email already exists")
		default:
			shared.WriteError(ctx, http.StatusBadGateway, "Something bad happened")
		}
		return
	}

	shared.WriteSuccess(ctx, http.StatusCreated, gin.H{"user": userResponse})
}

func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var payload *models.SignInInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		shared.WriteError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, refreshToken, err := ac.authService.SignIn(payload)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			shared.WriteError(ctx, http.StatusBadRequest, "Invalid email or Password")
		} else {
			shared.WriteError(ctx, http.StatusBadGateway, "Something bad happened")
		}
		return
	}

	ctx.SetCookie("access_token", accessToken, ac.accessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, ac.refreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", ac.accessTokenMaxAge*60, "/", "localhost", false, false)

	shared.WriteSuccess(ctx, http.StatusOK, gin.H{"access_token": accessToken})
}

// Refresh Access Token
func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {
	cookie, err := ctx.Cookie("refresh_token")

	if err != nil {
		shared.WriteError(ctx, http.StatusForbidden, "could not refresh access token")
		return
	}

	accessToken, err := ac.authService.RefreshAccessToken(cookie)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrRefreshTokenInvalid):
			shared.WriteError(ctx, http.StatusForbidden, "could not refresh access token")
		case errors.Is(err, services.ErrRefreshUserNotFound):
			shared.WriteError(ctx, http.StatusForbidden, err.Error())
		default:
			shared.WriteError(ctx, http.StatusBadGateway, "Something bad happened")
		}
		return
	}

	ctx.SetCookie("access_token", accessToken, ac.accessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", ac.accessTokenMaxAge*60, "/", "localhost", false, false)

	shared.WriteSuccess(ctx, http.StatusOK, gin.H{"access_token": accessToken})
}

func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, false)

	shared.WriteSuccess(ctx, http.StatusOK, gin.H{"message": "logout successful"})
}
