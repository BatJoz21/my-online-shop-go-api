package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateRefreshToken(userId int64, email, role string) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    userId,
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	return refreshToken.SignedString([]byte(secretkey))
}
