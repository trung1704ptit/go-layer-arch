package shared

import "github.com/gin-gonic/gin"

func WriteError(ctx *gin.Context, statusCode int, message string) {
	ctx.AbortWithStatusJSON(statusCode, gin.H{"error": message})
}

func WriteSuccess(ctx *gin.Context, statusCode int, data interface{}) {
	ctx.JSON(statusCode, gin.H{"data": data})
}
