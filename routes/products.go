package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/BatJoz21/my-online-shop-go-api/models"
	"github.com/BatJoz21/my-online-shop-go-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func createNewProduct(context *gin.Context) {
	category_id, err := strconv.ParseInt(context.PostForm("category_id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	price, err := decimal.NewFromString(context.PostForm("price"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var image *string
	file, err := context.FormFile("image")
	if err == nil {
		image, err = utils.SaveProductImage(file, context)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	product := models.Product{
		CategoryID:  category_id,
		Name:        context.PostForm("name"),
		Slug:        context.PostForm("slug"),
		Description: context.PostForm("description"),
		Price:       price,
		Image:       image,
	}
	err = product.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Product created", "product": product})
}

func getAllProducts(context *gin.Context) {
	category := context.DefaultQuery("category", "")

	page, err := strconv.Atoi(context.DefaultQuery("page", "1"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	offset := 0 + (models.ProductPerPageLimit * (page - 1))

	products, err := models.GetAllProducts(category, offset)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, products)
}

func getProductImage(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	product, err := models.GetActiveProduct(id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	if product.Image == nil {
		context.Status(http.StatusNoContent)
		return
	}

	path := utils.GetProductImagePath(product.Image)

	context.File(path)
}

func getProduct(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	product, err := models.GetActiveProduct(id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, product)
}

func updateProduct(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"is_updated": false,
			"message": err.Error()})
		return
	}

	product, err := models.GetActiveProduct(id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"is_updated": false,
			"message": err.Error()})
		return
	}

	category_id, err := strconv.ParseInt(context.PostForm("category_id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"is_updated": false,
			"message": err.Error()})
		return
	}
	price, err := decimal.NewFromString(context.PostForm("price"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"is_updated": false,
			"message": err.Error()})
		return
	}

	file, err := context.FormFile("image")
	if err == nil {
		err = utils.RemoveProductImage(product.Image)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"is_updated": false,
				"message": err.Error()})
			return
		}

		image, err := utils.SaveProductImage(file, context)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"is_updated": false,
				"message": err.Error()})
			return
		}

		product.Image = image
	}

	product.CategoryID = category_id
	product.Name = context.PostForm("name")
	product.Slug = context.PostForm("slug")
	product.Description = context.PostForm("description")
	product.Price = price

	err = product.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"is_updated": false,
			"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"is_updated": true, "message": "Product updated"})
}

func restoreSoftDeletedProduct(context *gin.Context) {
	// Get existing soft deleted product
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"is_restored": false, "message": err.Error()})
		return
	}
	product, err := models.GetProduct(id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"is_restored": false, "message": err.Error()})
		return
	}

	// Restoring product
	err = product.Restore()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"is_restored": false,
			"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"is_restored": true, "message": "Product is active"})
}

func softDeleteProduct(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"is_softDelete": false, "message": err.Error()})
		return
	}

	product, err := models.GetActiveProduct(id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"is_softDelete": false, "message": err.Error()})
		return
	}
	fmt.Println(product.ID)

	err = product.SoftDelete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"is_softDelete": false,
			"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"is_softDelete": true, "message": "Product not active"})
}

func deleteProduct(context *gin.Context) {
	// Get existing product
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"is_delete": false, "message": err.Error()})
		return
	}
	product, err := models.GetProduct(id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"is_delete": false, "message": err.Error()})
		return
	}

	// Delete all variant
	err = models.DeleteAllVariantOfAProduct(product.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"is_delete": false,
			"message": err.Error()})
		return
	}

	// Delete product image
	err = utils.RemoveProductImage(product.Image)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"is_delete": false,
			"message": err.Error()})
		return
	}

	// Delete product
	err = product.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"is_delete": false,
			"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"is_delete": true, "message": "Product deleted"})
}
