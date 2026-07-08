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
	custGroup.GET("products/:id/variants", getAllProductVariants)

	merchantGroup := server.Group("/")
	merchantGroup.Use(middlewares.Authenticate)
	merchantGroup.Use(middlewares.MerchantMiddleware())
	merchantGroup.POST("products", createNewProduct)
	merchantGroup.PUT("products/:id", updateProduct)
	merchantGroup.PUT("products/:id/restore", restoreSoftDeletedProduct)
	merchantGroup.DELETE("products/:id", softDeleteProduct)
	merchantGroup.DELETE("products/:id/delete", deleteProduct)
	merchantGroup.POST("products/:id/variants", createProductVariant)
	merchantGroup.PUT("products/:id/variants/:variant_id", updateProductVariant)
	merchantGroup.PUT("products/:id/variants/:variant_id/stock", updateVariantStock)
	merchantGroup.DELETE("products/:id/variants/:variant_id", deleteVariant)
}
