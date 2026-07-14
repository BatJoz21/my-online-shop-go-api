package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/BatJoz21/my-online-shop-go-api/models"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	_ "github.com/shopspring/decimal"
)

func generateNewOrder(context *gin.Context) {
	var dto models.NewOrderDTO
	err := context.ShouldBindJSON(&dto)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	order := models.Order{
		UserID:          context.GetInt64("userId"),
		OrderNumber:     "",
		Status:          "pending",
		ShippingAddress: dto.ShippingAddress,
	}

	err = order.GenerateNew()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, order)
}

func getAllUsersOrders(context *gin.Context) {
	uID := context.GetInt64("userId")

	orders, err := models.GetOrders(uID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, orders)
}

func populateGeneratedOrder(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("orderID"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	o, err := models.GetOrder(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Generate order number
	o.OrderNumber = fmt.Sprintf("ORD-%s-%06d", time.Now().Format("20060102"), o.ID)

	// Get total price
	totalAmount, err := models.GetTotalAmountFromOrderItems(o.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	o.TotalAmount, err = decimal.NewFromString(*totalAmount)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	err = o.InputData()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Finished setting up Order"})
}

func showOrderDetail(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("orderID"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	order, orderItems, err := models.GetOrderForShowPage(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"order": order, "orderItems": orderItems})
}

func deleteOrder(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("orderID"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	o, err := models.GetOrder(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	err = o.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
