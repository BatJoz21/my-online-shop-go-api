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

func VerifyAccessToken(token string) error {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Unexpected signing method")
		}

		return secretkey, nil
	})
	if err != nil {
		return err
	}

	if !parsedToken.Valid {
		return errors.New("Invalid token")
	}

	_, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("Invalid token claims")
	}

	// user_id := claims["user_id"].(int)
	// email := claims["email"].(string)
	// role := claims["role"].(string)

	return nil
}
