package middleware

import (
	"api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Token() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("access_token")
		if err != nil {
			ctx.Next()
			return
		}

		user, err := utils.DecryptToken(token)
		if err != nil {
			ctx.Next()
			return
		}

		ctx.Set("UserId", user.ID)
		ctx.Next()
	}
}

func AuthGuard() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, exists := ctx.Get("UserId")
		if !exists || id == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: You must be logged in to access this resource."})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
