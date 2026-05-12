package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-layer-arch/internal/initializers"
)

func RequireAPIKey() gin.HandlerFunc {
	config, _ := initializers.LoadConfig()
	expectedKey := config.BackendAPIKey

	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("X-API-Key")
		if expectedKey == "" || subtle.ConstantTimeCompare([]byte(apiKey), []byte(expectedKey)) != 1 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "invalid api key"})
			return
		}

		ctx.Next()
	}
}
