package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var key = os.Getenv("TOKEN_KEY")

func GenerateToken(user_id int64, email, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user_id,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Minute * 30).Unix(),
	})

	return token.SignedString([]byte(key))
}
