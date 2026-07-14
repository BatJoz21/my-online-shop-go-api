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

	server.GET("products", getAllStockedProducts)
	server.GET("products/:id/image", getProductImage)
	server.GET("categories", getCategories)

	custGroup := server.Group("/")
	custGroup.Use(middlewares.Authenticate)
	custGroup.GET("products/:id", getStockedProduct)
	custGroup.GET("products/:id/variants", getAllProductVariants)

	custGroup.GET("cart", getCartID)
	custGroup.GET("cart/total", getTotalItemOnCart)
	custGroup.POST("addToCart", addItemToCart)
	custGroup.GET("carts", getAllItemOnCart)
	custGroup.PUT("carts/:id", updateItemInCart)
	custGroup.DELETE("carts/:id", removeItemFromCart)
	custGroup.DELETE("carts", removeAllItemFromCart)

	custGroup.GET("orders", getAllUsersOrders)
	custGroup.POST("orders", generateNewOrder)
	custGroup.GET("orders/:orderID", showOrderDetail)
	custGroup.PUT("orders/:orderID", populateGeneratedOrder)
	custGroup.DELETE("orders/:orderID", deleteOrder)
	custGroup.POST("orders/:orderID/items", addItemToOrder)
	custGroup.GET("orders/:orderID/items", getAllItemsFromAnOrder)
	custGroup.DELETE("orders/:orderID/items/:orderItemID", deleteOrderItem)

	merchantGroup := server.Group("/")
	merchantGroup.Use(middlewares.Authenticate)
	merchantGroup.Use(middlewares.MerchantMiddleware())
	merchantGroup.GET("products/all", getAllProducts)
	merchantGroup.GET("products/all/:id", getProduct)
	merchantGroup.POST("products", createNewProduct)
	merchantGroup.PUT("products/:id", updateProduct)
	merchantGroup.PUT("products/:id/restore", restoreSoftDeletedProduct)
	merchantGroup.DELETE("products/:id", softDeleteProduct)
	merchantGroup.DELETE("products/:id/delete", deleteProduct)

	merchantGroup.POST("products/:id/variants", createProductVariant)
	merchantGroup.GET("products/:id/variants/:variant_id", getProductVariant)
	merchantGroup.PUT("products/:id/variants/:variant_id", updateProductVariant)
	merchantGroup.PUT("products/:id/variants/:variant_id/stock", updateVariantStock)
	merchantGroup.DELETE("products/:id/variants/:variant_id", deleteVariant)
}
