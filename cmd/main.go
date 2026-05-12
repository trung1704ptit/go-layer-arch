package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-layer-arch/internal/controllers"
	"github.com/go-layer-arch/internal/initializers"
	"github.com/go-layer-arch/internal/repositories"
	"github.com/go-layer-arch/internal/routes"
	"github.com/go-layer-arch/internal/services"
	"github.com/go-layer-arch/pkg/shared"
)

var (
	server              *gin.Engine
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	PostController      controllers.PostController
	PostRouteController routes.PostRouteController
)

func init() {
	config, err := initializers.LoadConfig()
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	authRepository := repositories.NewAuthRepository(initializers.DB)
	authService := services.NewAuthService(config, authRepository)
	AuthController = controllers.NewAuthController(authService, config.AccessTokenMaxAge, config.RefreshTokenMaxAge)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	userService := services.NewUserService()
	UserController = controllers.NewUserController(userService)
	UserRouteController = routes.NewRouteUserController(UserController)

	postRepository := repositories.NewPostRepository(initializers.DB)
	postService := services.NewPostService(postRepository)
	PostController = controllers.NewPostController(postService)
	PostRouteController = routes.NewRoutePostController(PostController)

	server = gin.Default()
}

func main() {
	config, err := initializers.LoadConfig()
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", config.ClientOrigin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Welcome to Golang with Gorm and Postgres"
		shared.WriteSuccess(ctx, http.StatusOK, gin.H{"message": message})
	})

	AuthRouteController.AuthRoute(router)
	UserRouteController.UserRoute(router)
	PostRouteController.PostRoute(router)
	log.Fatal(server.Run(":" + config.ServerPort))
}
