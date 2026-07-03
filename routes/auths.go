package routes

import (
	"net/http"

	"github.com/BatJoz21/my-online-shop-go-api/models"
	"github.com/BatJoz21/my-online-shop-go-api/utils"
	"github.com/gin-gonic/gin"
)

func signup(context *gin.Context) {
	// Get user's input
	var signupData models.UserSignUp
	err := context.ShouldBindJSON(&signupData)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Hashing password
	hashedPassword, err := utils.HashPassword(signupData.Password)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Save new user to Database
	newUser := models.User{
		Name:         signupData.Name,
		Email:        signupData.Email,
		PasswordHash: hashedPassword,
		Role:         models.RoleCustomer,
	}
	err = newUser.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "New user created"})
}

func login(context *gin.Context) {
	// Get user's input
	var loginData models.UserLogin
	err := context.ShouldBindJSON(&loginData)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Validate user's input
	err = models.ValidateCredentials(&loginData)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := models.GetUserDataForSession(loginData.Email)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Generate JWT token
	accessToken, err := utils.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login successful", "jwt_token": accessToken, "user": user})
}
