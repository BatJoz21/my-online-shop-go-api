package routes

import (
	"net/http"
	"strconv"

	"github.com/BatJoz21/my-online-shop-go-api/models"
	"github.com/gin-gonic/gin"
)

func createProductVariant(context *gin.Context) {
	var newVariant models.ProductVariantDTO
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
		Stock:         newVariant.Stock,
	}

	err = variant.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "New variant created", "product_variant": variant})
}

func getAllProductVariants(context *gin.Context) {
	product_id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	variants, err := models.GetAllVariantOfAProduct(product_id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, variants)
}

func updateProductVariant(context *gin.Context) {
	// Get existing product
	variant := getExistingVariant(context)

	// Get input data
	var variantDTO models.ProductVariantDTO
	err := context.ShouldBindJSON(&variantDTO)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Set input value to struct
	variant.Name = variantDTO.Name
	variant.Sku = variantDTO.Sku
	variant.PriceModifier = variantDTO.PriceModifier
	variant.Stock = variantDTO.Stock

	// Update
	err = variant.Update()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"is_updated": true, "message": "Product variant updated"})
}

func updateVariantStock(context *gin.Context) {
	// Get existing variant
	variant := getExistingVariant(context)

	var dto models.UpdateStockDTO
	err := context.ShouldBindJSON(&dto)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	variant.Stock += dto.Stock
	err = variant.UpdateStock()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"stock": variant.Stock})
}

func deleteVariant(context *gin.Context) {
	// Get existing variant
	variant := getExistingVariant(context)

	err := variant.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Product variant deleted"})
}

func getExistingVariant(context *gin.Context) *models.ProductVariant {
	id, err := strconv.ParseInt(context.Param("variant_id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return nil
	}
	variant, err := models.GetVariant(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return nil
	}

	return variant
}
