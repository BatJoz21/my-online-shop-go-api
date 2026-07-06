package middlewares

import (
	"net/http"
	"strings"

	"github.com/BatJoz21/my-online-shop-go-api/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	// Get logged in user token
	token := context.Request.Header.Get("Authorization")
	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
		return
	}
	token = strings.TrimPrefix(token, "Bearer ")

	// Verify token
	userId, email, role, err := utils.VerifyAccessToken(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// Set variable and value on context
	context.Set("userId", userId)
	context.Set("userEmail", email)
	context.Set("userRole", role)

	context.Next()
}
