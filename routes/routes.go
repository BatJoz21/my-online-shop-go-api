package routes

import (
	"github.com/BatJoz21/my-online-shop-go-api/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.POST("register", signup)
	server.POST("login", login)
	server.POST("logout", logout)
	server.POST("refresh", refreshJWT)

	custGroup := server.Group("/")
	custGroup.Use(middlewares.Authenticate)
}
