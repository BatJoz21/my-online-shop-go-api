package routes

import (
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
	stock, err := strconv.ParseInt(context.PostForm("stock"), 10, 64)
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
		Stock:       stock,
		Image:       image,
	}
	err = product.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Product created", "product": product})
}
