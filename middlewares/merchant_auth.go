package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MerchantMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		role := context.GetString("userRole")

		if role != "merchant" && role != "admin" {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "access denied"})
			return
		}

		context.Next()
	}
}
