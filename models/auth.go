package models

import (
	"errors"

	"github.com/BatJoz21/my-online-shop-go-api/database"
	"github.com/BatJoz21/my-online-shop-go-api/utils"
)

func ValidateCredentials(u *UserLogin) error {
	query := `SELECT password_hash FROM users WHERE email = ?`
	row := database.DB.QueryRow(query, u.Email)

	var retrievedPassword string
	err := row.Scan(&retrievedPassword)
	if err != nil {
		return err
	}

	if !utils.CheckPasswordHash(u.Password, retrievedPassword) {
		return errors.New("Invalid credentials")
	}

	return nil
}
