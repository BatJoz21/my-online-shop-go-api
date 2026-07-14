package routes

import (
	"net/http"
	"strconv"

	"github.com/BatJoz21/my-online-shop-go-api/models"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func addItemToCart(context *gin.Context) {
	cartID, err := models.GetUserCartID(context.GetInt64("userId"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var dto models.CartItemsDTO
	err = context.ShouldBindJSON(&dto)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	priceSnapshot, err := decimal.NewFromString(dto.PriceSnapshot)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	cartItem := models.CartItems{
		CartID:        *cartID,
		ProductID:     dto.ProductID,
		VariantID:     dto.VariantID,
		Quantity:      dto.Quantity,
		PriceSnapshot: priceSnapshot,
	}
	err = cartItem.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Item has been added to your cart"})
}

func getAllItemOnCart(context *gin.Context) {
	cartID, err := models.GetUserCartID(context.GetInt64("userId"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	items, err := models.GetAllItemInCart(*cartID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, items)
}

func getTotalItemOnCart(context *gin.Context) {
	cartID, err := models.GetUserCartID(context.GetInt64("userId"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	total, err := models.GetTotalItemInACart(*cartID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, total)
}

func updateItemInCart(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var dto models.UpdateCartItemDTO
	err = context.ShouldBindJSON(&dto)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	priceSnapshot, err := decimal.NewFromString(dto.PriceSnapshot)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	item := models.CartItems{
		ID:            id,
		VariantID:     dto.VariantID,
		Quantity:      dto.Quantity,
		PriceSnapshot: priceSnapshot,
	}

	err = item.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Cart item updated"})
}

func removeItemFromCart(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = models.DeleteCartItem(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
}

func removeAllItemFromCart(context *gin.Context) {
	cartID, err := models.GetUserCartID(context.GetInt64("userId"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	err = models.EmptyCart(*cartID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Your cart is now empty"})
}
