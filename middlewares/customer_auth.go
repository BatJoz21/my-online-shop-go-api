package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CustomerMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		role := context.GetString("userRole")

		if role != "customer" && role != "superadmin" {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "access denied"})
			return
		}

		context.Next()
	}
}
