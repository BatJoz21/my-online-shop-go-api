package routes

import (
	"net/http"
	"strconv"

	"github.com/BatJoz21/my-online-shop-go-api/models"
	"github.com/gin-gonic/gin"
)

func getAllUsers(context *gin.Context) {
	page, err := strconv.Atoi(context.DefaultQuery("page", "1"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	search := context.DefaultQuery("search", "")
	role := context.DefaultQuery("role", "")

	offset := models.UserPerPageLimit * (page - 1)

	users, err := models.GetUsers(search, role, offset)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, users)
}

func getUser(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("uID"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := models.GetUser(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, user)
}

func getUserProfile(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("uID"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := models.GetUser(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, user)
}

func updateUserRole(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("uID"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var dto models.UpdateRoleDTO
	if err := context.ShouldBindJSON(&dto); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if err := models.UpdateRole(id, dto.Role); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Role updated"})
}
