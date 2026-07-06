package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretkey = os.Getenv("TOKEN_KEY")

func GenerateToken(user_id int64, email, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user_id,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Minute * 30).Unix(),
	})

	return token.SignedString([]byte(secretkey))
}

func VerifyAccessToken(token string) (int64, string, string, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Unexpected signing method")
		}

		return []byte(secretkey), nil
	})
	if err != nil {
		return 0, "", "", err
	}

	if !parsedToken.Valid {
		return 0, "", "", errors.New("Invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", "", errors.New("Invalid token claims")
	}

	user_id := int64(claims["user_id"].(float64))
	email := claims["email"].(string)
	role := claims["role"].(string)

	return user_id, email, role, nil
}
