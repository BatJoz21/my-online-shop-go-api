package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuperAdminMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		role := context.GetString("userRole")

		if role != "superadmin" {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Access denied"})
			return
		}

		context.Next()
	}
}
