package routes

import (
	"net/http"
	"strconv"

	"github.com/BatJoz21/my-online-shop-go-api/models"
	"github.com/gin-gonic/gin"
)

func createProductVariant(context *gin.Context) {
	var newVariant models.NewProductVariant
	err := context.ShouldBindJSON(&newVariant)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	productId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	variant := models.ProductVariant{
		ProductID:     productId,
		Name:          newVariant.Name,
		Sku:           newVariant.Sku,
		PriceModifier: newVariant.PriceModifier,
	}

	err = variant.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "New variant created", "product_variant": variant})
}
