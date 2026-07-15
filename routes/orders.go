package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/BatJoz21/my-online-shop-go-api/models"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

const (
	timeLayout      = "2006-01-02"
	shippingCostStr = "50000.00"
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

func getAllOrder(context *gin.Context) {
	status := context.DefaultQuery("status", "")

	orders, err := models.GetAllOrders(status)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, orders)
}

func getAllUsersOrders(context *gin.Context) {
	uID := context.GetInt64("userId")
	status := context.DefaultQuery("status", "")

	orders, err := models.GetOrders(uID, status)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, orders)
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

func showOrderDetailForMerchant(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("orderID"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	order, orderItems, err := models.GetOrderForMerchantShowPage(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"order": order, "orderItems": orderItems})
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
	shippingCost, err := decimal.NewFromString(shippingCostStr)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	o.TotalAmount = o.TotalAmount.Add(shippingCost)

	err = o.InputData()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Finished setting up Order"})
}

func editOrder(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("orderID"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	order, err := models.GetOrder(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var editOrder models.EditOrderDTO
	err = context.ShouldBindJSON(&editOrder)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if editOrder.EstimatedArrival != "" {
		parsedTime, err := time.Parse(timeLayout, editOrder.EstimatedArrival)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		order.EstimatedArrival = &parsedTime
	}

	order.Status = editOrder.Status
	order.ShippingAddress = editOrder.ShippingAddress

	err = order.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Order updated"})
}

func completeOrder(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("orderID"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	order, err := models.GetOrder(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var dto models.ChangeStatusOrderDTO
	err = context.ShouldBindJSON(&dto)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	order.Status = dto.Status

	err = order.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Your order has completed. Thank you."})
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
