package main

import (
	"flag"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	controllerMarker = "\t// FEATURE CONTROLLERS END"
	wiringMarker     = "\t// FEATURE WIRING END"
	routeMarker      = "\t// FEATURE ROUTES END"
)

type feature struct {
	Snake     string
	Pascal    string
	Camel     string
	RoutePath string
}

type generatedFile struct {
	path    string
	content string
}

func main() {
	name := flag.String("name", "", "feature name, for example my_feature")
	flag.Parse()

	feat, err := parseFeatureName(*name)
	if err != nil {
		exit(err)
	}

	if err := createFeatureFiles(feat); err != nil {
		exit(err)
	}

	if err := updateMain(feat); err != nil {
		exit(err)
	}

	fmt.Printf("Created feature scaffold for %q\n", feat.Snake)
}

func parseFeatureName(raw string) (feature, error) {
	name := strings.TrimSpace(raw)
	if name == "" {
		return feature{}, fmt.Errorf("feature name is required, usage: make feature name=my_feature")
	}

	normalized := strings.NewReplacer("-", "_", " ", "_").Replace(name)
	if !regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*$`).MatchString(normalized) {
		return feature{}, fmt.Errorf("feature name must start with a letter and contain only letters, numbers, underscores, hyphens, or spaces")
	}

	parts := strings.FieldsFunc(normalized, func(r rune) bool {
		return r == '_'
	})
	for i, part := range parts {
		if part == "" {
			return feature{}, fmt.Errorf("feature name contains an empty word")
		}
		parts[i] = strings.ToLower(part)
	}

	pascalParts := make([]string, 0, len(parts))
	for _, part := range parts {
		pascalParts = append(pascalParts, strings.ToUpper(part[:1])+part[1:])
	}

	pascal := strings.Join(pascalParts, "")
	return feature{
		Snake:     strings.Join(parts, "_"),
		Pascal:    pascal,
		Camel:     strings.ToLower(pascal[:1]) + pascal[1:],
		RoutePath: pluralize(strings.Join(parts, "-")),
	}, nil
}

func pluralize(word string) string {
	if strings.HasSuffix(word, "s") {
		return word
	}

	if strings.HasSuffix(word, "y") && len(word) > 1 {
		previous := word[len(word)-2]
		if !strings.ContainsRune("aeiou", rune(previous)) {
			return word[:len(word)-1] + "ies"
		}
	}

	return word + "s"
}

func createFeatureFiles(feat feature) error {
	files := []generatedFile{
		{
			path:    filepath.Join("internal", "models", feat.Snake+".model.go"),
			content: modelTemplate(feat),
		},
		{
			path:    filepath.Join("internal", "repositories", feat.Snake+".repository.go"),
			content: repositoryTemplate(feat),
		},
		{
			path:    filepath.Join("internal", "services", feat.Snake+".service.go"),
			content: serviceTemplate(feat),
		},
		{
			path:    filepath.Join("internal", "controllers", feat.Snake+".controller.go"),
			content: controllerTemplate(feat),
		},
		{
			path:    filepath.Join("internal", "routes", feat.Snake+".routes.go"),
			content: routeTemplate(feat),
		},
	}

	for _, file := range files {
		if err := writeGoFile(file.path, file.content); err != nil {
			return err
		}
	}

	return nil
}

func writeGoFile(path string, content string) error {
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("%s already exists", path)
	} else if !os.IsNotExist(err) {
		return err
	}

	formatted, err := format.Source([]byte(content))
	if err != nil {
		return fmt.Errorf("format %s: %w", path, err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	return os.WriteFile(path, formatted, 0644)
}

func updateMain(feat feature) error {
	path := filepath.Join("cmd", "main.go")
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	content := string(contentBytes)
	content, err = insertBeforeMarker(content, controllerMarker, feat.Pascal+"Controller", fmt.Sprintf(
		"\t%sController      controllers.%sController\n\t%sRouteController routes.%sRouteController\n",
		feat.Pascal,
		feat.Pascal,
		feat.Pascal,
		feat.Pascal,
	))
	if err != nil {
		return err
	}

	content, err = insertBeforeMarker(content, wiringMarker, "New"+feat.Pascal+"Repository", fmt.Sprintf(
		"\t%sRepository := repositories.New%sRepository(initializers.DB)\n\t%sService := services.New%sService(%sRepository)\n\t%sController = controllers.New%sController(%sService)\n\t%sRouteController = routes.NewRoute%sController(%sController)\n\n",
		feat.Camel,
		feat.Pascal,
		feat.Camel,
		feat.Pascal,
		feat.Camel,
		feat.Pascal,
		feat.Pascal,
		feat.Camel,
		feat.Pascal,
		feat.Pascal,
		feat.Pascal,
	))
	if err != nil {
		return err
	}

	content, err = insertBeforeMarker(content, routeMarker, feat.Pascal+"Route(router)", fmt.Sprintf(
		"\t%sRouteController.%sRoute(router)\n",
		feat.Pascal,
		feat.Pascal,
	))
	if err != nil {
		return err
	}

	formatted, err := format.Source([]byte(content))
	if err != nil {
		return fmt.Errorf("format %s: %w", path, err)
	}

	return os.WriteFile(path, formatted, 0644)
}

func insertBeforeMarker(content string, marker string, unique string, addition string) (string, error) {
	if strings.Contains(content, unique) {
		return content, nil
	}

	index := strings.Index(content, marker)
	if index == -1 {
		return "", fmt.Errorf("missing scaffold marker %q in cmd/main.go", strings.TrimSpace(marker))
	}

	return content[:index] + addition + content[index:], nil
}

func modelTemplate(feat feature) string {
	return fmt.Sprintf(`package models

type %[1]s struct {
	// Add model fields here.
}

type Create%[1]sRequest struct {
	// Add request fields here.
}

type Update%[1]sRequest struct {
	// Add request fields here.
}
`, feat.Pascal)
}

func repositoryTemplate(feat feature) string {
	return fmt.Sprintf(`package repositories

import "gorm.io/gorm"

type %[1]sRepository interface {
	// Add repository methods here.
}

type %[2]sRepository struct {
	db *gorm.DB
}

func New%[1]sRepository(db *gorm.DB) %[1]sRepository {
	return &%[2]sRepository{db: db}
}
`, feat.Pascal, feat.Camel)
}

func serviceTemplate(feat feature) string {
	return fmt.Sprintf(`package services

import "github.com/go-layer-arch/internal/repositories"

type %[1]sService interface {
	// Add service methods here.
}

type %[2]sService struct {
	%[2]sRepository repositories.%[1]sRepository
}

func New%[1]sService(%[2]sRepository repositories.%[1]sRepository) %[1]sService {
	return &%[2]sService{%[2]sRepository: %[2]sRepository}
}
`, feat.Pascal, feat.Camel)
}

func controllerTemplate(feat feature) string {
	return fmt.Sprintf(`package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-layer-arch/internal/services"
	"github.com/go-layer-arch/pkg/shared"
)

type %[1]sController struct {
	%[2]sService services.%[1]sService
}

func New%[1]sController(%[2]sService services.%[1]sService) %[1]sController {
	return %[1]sController{%[2]sService: %[2]sService}
}

func (fc *%[1]sController) Placeholder(ctx *gin.Context) {
	shared.WriteSuccess(ctx, http.StatusOK, gin.H{"message": "%[3]s feature placeholder"})
}
`, feat.Pascal, feat.Camel, feat.Snake)
}

func routeTemplate(feat feature) string {
	return fmt.Sprintf(`package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-layer-arch/internal/controllers"
)

type %[1]sRouteController struct {
	%[2]sController controllers.%[1]sController
}

func NewRoute%[1]sController(%[2]sController controllers.%[1]sController) %[1]sRouteController {
	return %[1]sRouteController{%[2]sController: %[2]sController}
}

func (rc *%[1]sRouteController) %[1]sRoute(rg *gin.RouterGroup) {
	router := rg.Group("%[3]s")
	router.GET("/", rc.%[2]sController.Placeholder)
}
`, feat.Pascal, feat.Camel, feat.RoutePath)
}

func exit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
