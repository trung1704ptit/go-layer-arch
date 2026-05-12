package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateFeatureScaffold(t *testing.T) {
	workingDir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatal(err)
		}
	})

	if err := os.Chdir(workingDir); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll("cmd", 0755); err != nil {
		t.Fatal(err)
	}

	mainFile := `package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-layer-arch/internal/controllers"
	"github.com/go-layer-arch/internal/initializers"
	"github.com/go-layer-arch/internal/repositories"
	"github.com/go-layer-arch/internal/routes"
	"github.com/go-layer-arch/internal/services"
)

var (
	server *gin.Engine

	PostController      controllers.PostController
	PostRouteController routes.PostRouteController
)

func init() {
	postRepository := repositories.NewPostRepository(initializers.DB)
	postService := services.NewPostService(postRepository)
	PostController = controllers.NewPostController(postService)
	PostRouteController = routes.NewRoutePostController(PostController)

	server = gin.Default()
}

func main() {
	router := server.Group("/api")
	PostRouteController.PostRoute(router)

	log.Fatal(server.Run(":8000"))
}
`

	if err := os.WriteFile(filepath.Join("cmd", "main.go"), []byte(mainFile), 0644); err != nil {
		t.Fatal(err)
	}

	feat, err := parseFeatureName("my_feature")
	if err != nil {
		t.Fatal(err)
	}

	if err := createFeatureFiles(feat); err != nil {
		t.Fatal(err)
	}

	if err := updateMain(feat); err != nil {
		t.Fatal(err)
	}

	if err := createFeatureFiles(feat); err != nil {
		t.Fatal(err)
	}

	if err := updateMain(feat); err != nil {
		t.Fatal(err)
	}

	expectedFiles := []string{
		filepath.Join("internal", "models", "my_feature.model.go"),
		filepath.Join("internal", "repositories", "my_feature.repository.go"),
		filepath.Join("internal", "services", "my_feature.service.go"),
		filepath.Join("internal", "controllers", "my_feature.controller.go"),
		filepath.Join("internal", "routes", "my_feature.routes.go"),
	}
	for _, file := range expectedFiles {
		if _, err := os.Stat(file); err != nil {
			t.Fatalf("expected generated file %s: %v", file, err)
		}
	}

	updatedMain, err := os.ReadFile(filepath.Join("cmd", "main.go"))
	if err != nil {
		t.Fatal(err)
	}

	mainContent := string(updatedMain)
	for _, expected := range []string{
		"MyFeatureController",
		"controllers.MyFeatureController",
		"NewMyFeatureRepository",
		"MyFeatureRouteController.MyFeatureRoute(router)",
	} {
		if !strings.Contains(mainContent, expected) {
			t.Fatalf("expected cmd/main.go to contain %q", expected)
		}
	}
}
