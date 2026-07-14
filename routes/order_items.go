package routes

import (
	"net/http"
	"strconv"

	"github.com/BatJoz21/my-online-shop-go-api/models"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func addItemToOrder(context *gin.Context) {
	var dto models.NewOrderItemDTO
	err := context.ShouldBindJSON(&dto)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	productID, err := strconv.ParseInt(dto.ProductID, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	variantID, err := strconv.ParseInt(dto.VariantID, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	itemQuantity, err := strconv.Atoi(dto.Quantity)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	variant, err := models.GetVariant(variantID)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if variant.Stock < 1 || variant.Stock < int64(itemQuantity) {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Item is out of stock"})
		return
	}
	variant.Stock = variant.Stock - int64(itemQuantity)

	priceSnap, err := decimal.NewFromString(dto.PriceSnapshot)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	quantity := decimal.NewFromInt(int64(itemQuantity))
	subtotal := priceSnap.Mul(quantity)

	orderItem := models.OrderItem{
		OrderID:       dto.OrderID,
		ProductID:     productID,
		VariantID:     variantID,
		ProductName:   dto.ProductName,
		Quantity:      itemQuantity,
		PriceSnapshot: priceSnap,
		Subtotal:      subtotal,
	}

	err = orderItem.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	err = variant.UpdateStock()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Success", "orderItem": orderItem})
}

func getAllItemsFromAnOrder(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("orderID"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	orderItems, err := models.GetAllItemFromOrder(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, orderItems)
}

func deleteOrderItem(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("orderItemID"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	item, err := models.GetAnItemFromOrder(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	variant, err := models.GetVariant(item.VariantID)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	variant.Stock += int64(item.Quantity)

	err = item.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	err = variant.UpdateStock()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
