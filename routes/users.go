package routes

import (
	"net/http"

	"github.com/BatJoz21/my-online-shop-go-api/models"
	"github.com/gin-gonic/gin"
)

func getAllUsers(context *gin.Context) {
	users, err := models.GetUsers()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, users)
}
