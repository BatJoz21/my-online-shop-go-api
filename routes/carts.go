package routes

import (
	"net/http"

	"github.com/BatJoz21/my-online-shop-go-api/models"
	"github.com/gin-gonic/gin"
)

func getCartID(context *gin.Context) {
	uID := context.GetInt64("userId")

	id, err := models.GetUserCartID(uID)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Cart not found"})
		return
	}

	context.JSON(http.StatusOK, id);
}
