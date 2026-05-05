package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-layer-arch/controllers"
	"github.com/go-layer-arch/initializers"
	"github.com/go-layer-arch/routes"
	"github.com/go-layer-arch/services"
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
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	authService := services.NewAuthService(initializers.DB)
	AuthController = controllers.NewAuthController(authService)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	userService := services.NewUserService()
	UserController = controllers.NewUserController(userService)
	UserRouteController = routes.NewRouteUserController(UserController)

	postService := services.NewPostService(initializers.DB)
	PostController = controllers.NewPostController(postService)
	PostRouteController = routes.NewRoutePostController(PostController)

	server = gin.Default()
}

func main() {
	config, err := initializers.LoadConfig(".")
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
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	AuthRouteController.AuthRoute(router)
	UserRouteController.UserRoute(router)
	PostRouteController.PostRoute(router)
	log.Fatal(server.Run(":" + config.ServerPort))
}
