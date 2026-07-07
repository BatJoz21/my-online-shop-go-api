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

	server.GET("products", getAllProducts)
	server.GET("products/:id/image", getProductImage)
	server.GET("categories", getCategories)

	custGroup := server.Group("/")
	custGroup.Use(middlewares.Authenticate)
	custGroup.GET("products/:id", getProduct)

	merchantGroup := server.Group("/")
	merchantGroup.Use(middlewares.Authenticate)
	merchantGroup.Use(middlewares.MerchantMiddleware())
	merchantGroup.POST("products", createNewProduct)
	merchantGroup.POST("products/:id/variants", createProductVariant)
}
