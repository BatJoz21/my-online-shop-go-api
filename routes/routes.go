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
	server.GET("products/:id", getStockedProduct)

	server.GET("products/:id/variants", getAllProductVariants)

	server.POST("/payments/webhook", handlePaymentWebhook)

	server.GET("products/:id/reviews", getReviewsOfAProduct)

	authGroup := server.Group("")
	authGroup.Use(middlewares.Authenticate)
	authGroup.GET("/profile/:uID", getUserProfile)

	custGroup := authGroup.Group("/")
	custGroup.Use(middlewares.CustomerMiddleware())

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
	custGroup.PUT("orders/:orderID/populate", populateGeneratedOrder)
	custGroup.PUT("orders/:orderID/complete", completeOrder)
	custGroup.DELETE("orders/:orderID", deleteOrder)
	custGroup.POST("orders/:orderID/items", addItemToOrder)
	custGroup.GET("orders/:orderID/items", getAllItemsFromAnOrder)
	custGroup.DELETE("orders/:orderID/items/:orderItemID", deleteOrderItem)

	custGroup.POST("orders/:orderID/payment", initiatePayment)

	custGroup.POST("products/:id/reviews", addReview)

	merchantGroup := authGroup.Group("/merchant/")
	merchantGroup.Use(middlewares.MerchantMiddleware())
	merchantGroup.GET("dashboard/stats", getDashboardStatsData)
	merchantGroup.GET("dashboard/orders", getRecentOrdersData)
	merchantGroup.GET("dashboard/low-stocked", getLowStockProductsData)
	merchantGroup.GET("dashboard/review", getRecentReviewData)

	merchantGroup.POST("categories", addNewCategory)

	merchantGroup.POST("products", createNewProduct)
	merchantGroup.GET("products", getAllProducts)
	merchantGroup.GET("products/:id", getProduct)
	merchantGroup.PUT("products/:id", updateProduct)
	merchantGroup.PUT("products/:id/restore", restoreSoftDeletedProduct)
	merchantGroup.DELETE("products/:id", softDeleteProduct)
	merchantGroup.DELETE("products/:id/delete", deleteProduct)

	merchantGroup.POST("products/:id/variants", createProductVariant)
	merchantGroup.GET("products/:id/variants/:variant_id", getProductVariant)
	merchantGroup.PUT("products/:id/variants/:variant_id", updateProductVariant)
	merchantGroup.PUT("products/:id/variants/:variant_id/stock", updateVariantStock)
	merchantGroup.DELETE("products/:id/variants/:variant_id", deleteVariant)

	merchantGroup.GET("orders", getAllOrder)
	merchantGroup.GET("orders/:orderID", showOrderDetailForMerchant)
	merchantGroup.PUT("orders/:orderID", editOrder)

	superadminGroup := authGroup.Group("/admin/")
	superadminGroup.Use(middlewares.SuperAdminMiddleware())
	superadminGroup.GET("users", getAllUsers)
	superadminGroup.GET("users/:uID", getUser)
	superadminGroup.PUT("users/:uID/role", updateUserRole)
}
