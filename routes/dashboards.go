package routes

import (
	"net/http"

	"github.com/BatJoz21/my-online-shop-go-api/models"
	"github.com/gin-gonic/gin"
)

func getDashboardStatsData(context *gin.Context) {
	stats, err := models.GetDashboardStats()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch stats"})
		return
	}

	context.JSON(http.StatusOK, stats)
}

func getRecentOrdersData(context *gin.Context) {
	orders, err := models.GetRecentOrders()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch stats"})
		return
	}

	context.JSON(http.StatusOK, orders)
}

func getLowStockProductsData(context *gin.Context) {
	products, err := models.GetLowStockProducts()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch stats"})
		return
	}

	context.JSON(http.StatusOK, products)
}

func getRecentReviewData(context *gin.Context) {
	reviews, err := models.GetRecentReviews()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch stats"})
		return
	}

	context.JSON(http.StatusOK, reviews)
}
