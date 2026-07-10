package routes

import (
	"net/http"

	"github.com/BatJoz21/my-online-shop-go-api/models"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func generateNewOrder(context *gin.Context) {
	var dto models.NewOrderDTO
	err := context.ShouldBindJSON(&dto)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	totalAmount, err := decimal.NewFromString(dto.TotalAmount)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	order := models.Order{
		UserID: context.GetInt64("userId"),
		OrderNumber: "",
		Status: "pending",
		TotalAmount: totalAmount,
		ShippingAddress: dto.ShippingAddress,
	}

	err = order.GenerateNew()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, order)
}