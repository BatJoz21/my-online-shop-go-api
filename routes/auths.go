package routes

import (
	"net/http"
	"time"

	"github.com/BatJoz21/my-online-shop-go-api/models"
	"github.com/BatJoz21/my-online-shop-go-api/utils"
	"github.com/gin-gonic/gin"
)

func signup(context *gin.Context) {
	// Get user's input
	var signupData models.UserSignUp
	err := context.ShouldBindJSON(&signupData)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "isRegistered": false})
		return
	}

	// Hashing password
	hashedPassword, err := utils.HashPassword(signupData.Password)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "isRegistered": false})
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
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "isRegistered": false})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "New user created", "isRegistered": true})
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

	userData, err := models.GetUserDataForSession(loginData.Email)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Generate JWT token
	accessToken, err := utils.GenerateToken(userData.ID, userData.Email, string(userData.Role))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Generate Refresh Token
	refreshToken, err := utils.GenerateRefreshToken(userData.ID, userData.Email, string(userData.Role))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	expiresAt := time.Now().Add(time.Hour * 24 * 7)
	deviceName := "default"
	refreshTokenStruct := models.RefreshToken{
		UserID:     userData.ID,
		DeviceName: &deviceName,
		TokenHash:  utils.HashRefreshToken(refreshToken),
		ExpiresAt:  &expiresAt,
	}
	err = refreshTokenStruct.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message":       "Login successful",
		"jwt_token":     accessToken,
		"refresh_token": refreshToken,
		"user":          userData})
}

func logout(context *gin.Context) {
	var request models.RefreshTokenRequest
	err := context.ShouldBindJSON(&request)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "isLoggedOut": false})
		return
	}

	hashed := utils.HashRefreshToken(request.RefreshToken)
	refreshToken, err := models.GetRefreshTokenByHashedToken(hashed)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "isLoggedOut": false})
		return
	}

	err = refreshToken.Revoke()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "isLoggedOut": false})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "user logged out", "isLoggedOut": true})
}

func refreshJWT(context *gin.Context) {
	var request models.RefreshTokenRequest
	err := context.ShouldBindJSON(&request)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	hashed := utils.HashRefreshToken(request.RefreshToken)
	refreshTokenStruct, err := models.GetRefreshTokenByHashedToken(hashed)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if refreshTokenStruct.RevokedAt != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Token is revoked"})
		return
	}
	if time.Now().After(*refreshTokenStruct.ExpiresAt) {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Token is expired"})
		return
	}

	userData, err := models.GetUserDataForRefreshToken(refreshTokenStruct.UserID)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	jwtToken, err := utils.GenerateToken(userData.ID, userData.Email, string(userData.Role))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "session refreshed", "jwt_token": jwtToken})
}
